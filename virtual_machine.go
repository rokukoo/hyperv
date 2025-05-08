package hypervctl

import (
	"log"
	"os"

	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/memory"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/processor"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualMachineState = virtual_system.ComputerSystemState

const (
	StateRunning VirtualMachineState = virtual_system.Running
	StateStopped VirtualMachineState = virtual_system.Off
	StateSuspend VirtualMachineState = virtual_system.Saved
)

// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-computersystem
type VirtualMachine struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	SavePath       string `json:"save_path"`
	CpuCoreCount   int    `json:"cpu_core_count"`
	MemorySizeMB   int    `json:"memory_size"`
	computerSystem *virtual_system.ComputerSystem
}

// Start 启动虚拟机
func (vm *VirtualMachine) Start() error {
	return vm.computerSystem.Start()
}

// Stop 停止虚拟机
func (vm *VirtualMachine) Stop(force bool) error {
	return vm.computerSystem.Stop(force)
}

// ForceStop 强制停止虚拟机
func (vm *VirtualMachine) ForceStop() error {
	return vm.computerSystem.ForceStop()
}

// Shutdown 正常关闭虚拟机
func (vm *VirtualMachine) Shutdown() error {
	return vm.computerSystem.Stop(false)
}

// Reboot 重启虚拟机
func (vm *VirtualMachine) Reboot(force bool) error {
	return vm.computerSystem.Reboot(force)
}

// ForceReboot 强制重启虚拟机
func (vm *VirtualMachine) ForceReboot() error {
	return vm.computerSystem.ForceReboot()
}

// Resume 恢复虚拟机
func (vm *VirtualMachine) Resume() error {
	return vm.computerSystem.Resume()
}

// State 获取虚拟机实时状态
func (vm *VirtualMachine) State() VirtualMachineState {
	var err error
	var state VirtualMachineState

	if state, err = vm.computerSystem.GetState(); err != nil {
		log.Fatalf("failed to get vm state: %v", err)
	}

	return state
}

// Suspend 挂起虚拟机
// 由于 Hyper-V 平台原生挂起功能并非真正意义上的挂起, 而是保存虚拟机的状态, 因此这里的挂起操作实际上是保存虚拟机的状态
// 保存虚拟机的状态后, 可以通过 Resume 恢复虚拟机的运行
func (vm *VirtualMachine) Suspend() error {
	return vm.computerSystem.Save()
}

func (vm *VirtualMachine) update(cs *virtual_system.ComputerSystem) error {
	var err error
	var virtualSystemSettingData *virtual_system.VirtualSystemSettingData

	vm.computerSystem = cs

	if virtualSystemSettingData, err = vm.computerSystem.GetVirtualSystemSettingData(); err != nil {
		return err
	}

	vm.Name = cs.ElementName
	vm.Description = virtualSystemSettingData.Notes[0]
	vm.SavePath = virtualSystemSettingData.ConfigurationDataRoot

	vm.CpuCoreCount = int(vm.computerSystem.MustGetProcessorSettingData().VirtualQuantity)
	vm.MemorySizeMB = int(vm.computerSystem.MustGetMemorySettingData().VirtualQuantity)

	return nil
}

func NewVirtualMachine(cs *virtual_system.ComputerSystem) (*VirtualMachine, error) {
	var err error
	vm := &VirtualMachine{}
	if err = vm.update(cs); err != nil {
		return nil, err
	}
	return vm, nil
}

func (vm *VirtualMachine) Create() (err error) {
	var buildVM *VirtualMachine
	var builder *VirtualMachineBuilder

	if builder, err = NewVirtualMachineBuilder(); err != nil {
		return err
	}

	builder.PrepareSystemSettings(vm.Name, func(systemSettingsData *virtual_system.VirtualSystemSettingData) {
		systemSettingsData.ConfigurationDataRoot = vm.SavePath
		if vm.Description != "" {
			systemSettingsData.Notes = []string{vm.Description}
		}
	})

	builder.PrepareProcessorSettings(func(processorSettings *processor.ProcessorSettingData) {
		processorSettings.VirtualQuantity = uint64(vm.CpuCoreCount)
	})

	builder.PrepareMemorySettings(func(memorySettings *memory.MemorySettingsData) {
		memorySettings.VirtualQuantity = uint64(vm.MemorySizeMB)
		memorySettings.DynamicMemoryEnabled = false
	})

	if buildVM, err = builder.Build(); err != nil {
		if errors.Unwrap(err).(*wmiext.JobError).ErrorCode == 32769 {
			// 创建失败, 当前目录下已存在同名的虚拟机
			return ErrorVirtualMachineAlreadyExists
		}
		return err
	}

	if err = vm.update(buildVM.computerSystem); err != nil {
		return err
	}

	return err
}

// ModifySpecOptions 修改虚拟机规格选项
type ModifySpecOptions struct {
	cpuCoreCount int
	memorySizeMB int
	confirmStop  bool
}

type Option func(*ModifySpecOptions)

func WithCpuCoreCount(cpuCoreCount int) Option {
	return func(options *ModifySpecOptions) {
		options.cpuCoreCount = cpuCoreCount
	}
}

func WithMemorySize(memorySizeMB int) Option {
	return func(options *ModifySpecOptions) {
		options.memorySizeMB = memorySizeMB
	}
}

func WithStop(confirmStop bool) Option {
	return func(options *ModifySpecOptions) {
		options.confirmStop = confirmStop
	}
}

// Modify 修改虚拟机规格
func (vm *VirtualMachine) Modify(options ...Option) (ok bool, err error) {
	var originalState = vm.State()
	opts := new(ModifySpecOptions)
	for _, option := range options {
		option(opts)
	}
	if opts.confirmStop && vm.State() != StateStopped {
		if err = vm.computerSystem.ForceStop(); err != nil {
			return
		}
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	// 修改 CPU 核心数
	if opts.cpuCoreCount > 0 {
		// HyperV 要求如果虚拟机未关闭, 则不允许修改 CPU 核心数
		if vm.State() != StateStopped {
			return false, errors.New("vm must be stopped before modifying spec")
		}

		processorSettingData := vm.computerSystem.MustGetProcessorSettingData()
		processorSettingData.VirtualQuantity = uint64(opts.cpuCoreCount)

		if err = processorSettingData.Put("VirtualQuantity", processorSettingData.VirtualQuantity); err != nil {
			return false, err
		}

		if err = vmms.ModifyProcessorSettings(processorSettingData); err != nil {
			return false, err
		}

		vm.CpuCoreCount = opts.cpuCoreCount
	}
	if opts.memorySizeMB > 0 {
		// HyperV 允许虚拟机运行状态下, 修改内存大小
		memorySettingData := vm.computerSystem.MustGetMemorySettingData()

		memorySettingData.VirtualQuantity = uint64(opts.memorySizeMB)

		if err = memorySettingData.Put("VirtualQuantity", memorySettingData.VirtualQuantity); err != nil {
			return false, err
		}

		if err = vmms.ModifyMemorySettings(memorySettingData); err != nil {
			if errors.Unwrap(err).(*wmiext.JobError).ErrorCode == 32768 {
				// Error code 32768: The operation cannot be performed while the virtual machine is in its current state.
				// Maybe the virtual machine does not support resizing memory dynamically (with not stop), try to stop the virtual machine and try again
				if err = vm.computerSystem.ForceStop(); err != nil {
					return false, err
				}
				if err = vmms.ModifyMemorySettings(memorySettingData); err != nil {
					return false, err
				}
			} else {
				return false, err
			}
		}

		vm.MemorySizeMB = opts.memorySizeMB
	}

	if vm.State() != originalState {
		if err = vm.computerSystem.ChangeState(originalState); err != nil {
			return false, err
		}
	}

	return true, nil
}

// FindVirtualMachineByName 根据虚拟机名称获取虚拟机
// 
// 参数:
//   vmName: 虚拟机名称
// 返回:
//   []*VirtualMachine: 虚拟机列表
//   error: 错误
func FindVirtualMachineByName(vmName string) ([]*VirtualMachine, error) {
	service, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	vms, err := service.FindComputerSystemsByName(vmName)
	if err != nil {
		return nil, err
	}
	var virtualMachines []*VirtualMachine
	for _, vm := range vms {
		virtualMachine, err := NewVirtualMachine(vm)
		if err != nil {
			return nil, err
		}
		virtualMachines = append(virtualMachines, virtualMachine)
	}
	return virtualMachines, nil
}

// FirstVirtualMachineByName 根据虚拟机名称获取第一个虚拟机
func FirstVirtualMachineByName(vmName string) (*VirtualMachine, error) {
	vms, err := FindVirtualMachineByName(vmName)
	if err != nil {
		return nil, err
	}
	if len(vms) == 0 {
		return nil, wmiext.NotFound
	}
	return vms[0], nil
}

func MustFirstVirtualMachineByName(vmName string) *VirtualMachine {
	vm, err := FirstVirtualMachineByName(vmName)
	if err != nil {
		log.Fatalf("failed to find virtual machine: %v", err)
	}
	return vm
}

// CreateVirtualMachine 创建虚拟机
func CreateVirtualMachine(name string, savePath string, cpuCoreCount int, memorySize int) (*VirtualMachine, error) {
	var err error
	virtualMachine := &VirtualMachine{
		Name:         name,
		SavePath:     savePath,
		CpuCoreCount: cpuCoreCount,
		MemorySizeMB: memorySize,
	}
	if err = virtualMachine.Create(); err != nil {
		return nil, err
	}
	return virtualMachine, nil
}

// ListVirtualMachines 获取所有虚拟机
func ListVirtualMachines() (vms []*VirtualMachine, err error) {
	var vsms *virtual_system.VirtualSystemManagementService
	var vm *VirtualMachine
	vsms, err = virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	computerSystems, err := vsms.ListComputerSystems()
	if err != nil {
		return nil, err
	}
	for _, cs := range computerSystems {
		if vm, err = NewVirtualMachine(cs); err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	return
}

// DestroyVirtualMachineByName 根据名称删除虚拟机
func DestroyVirtualMachineByName(name string, del bool) (ok bool, err error) {
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	vm, err := FirstVirtualMachineByName(name)
	if err != nil {
		return false, err
	}
	if vm.State() != StateStopped {
		return false, errors.New("vm must be stopped before deleting")
	}

	if err = vmms.DestroySystem(vm.computerSystem); err != nil {
		return false, err
	}
	if del {
		// RemoveAll 可以删除非空文件夹
		if err = os.RemoveAll(vm.SavePath); err != nil {
			return false, err
		}
	}
	return true, nil
}

// DeleteVirtualMachineByName 根据名称删除虚拟机
func DeleteVirtualMachineByName(name string) (ok bool, err error) {
	return DestroyVirtualMachineByName(name, true)
}

// ModifySpec 修改虚拟机规格
func (vm *VirtualMachine) ModifySpec(cpuCoreCount, memorySize int) (ok bool, err error) {
	var options []Option
	if cpuCoreCount > 0 {
		vm.CpuCoreCount = cpuCoreCount
		options = append(options, WithCpuCoreCount(cpuCoreCount))
		options = append(options, WithStop(true))
	}
	if memorySize > 0 {
		vm.MemorySizeMB = memorySize
		options = append(options, WithMemorySize(memorySize))
	}
	return vm.Modify(options...)
}

// ModifyVirtualMachineSpecByName 根据虚拟机名称修改虚拟机规格
func ModifyVirtualMachineSpecByName(name string, cpuCoreCount int, memorySize int) (ok bool, err error) {
	vm, err := FirstVirtualMachineByName(name)
	if err != nil {
		return false, err
	}
	return vm.ModifySpec(cpuCoreCount, memorySize)
}
