package hypervctl

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/microsoft/wmi/pkg/virtualization/core/memory"
	"github.com/microsoft/wmi/pkg/virtualization/core/processor"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/wmiext"
	"log"
	"os"
	"strings"
)

type IVirtualMachine interface {
	Start() (ok bool, err error)
	Stop(force bool) (ok bool, err error)
	Reboot(force bool) (ok bool, err error)
	ForceStop() (ok bool, err error)
	ForceReboot() (ok bool, err error)
	Suspend() (ok bool, err error)
	Resume() (ok bool, err error)
	Create() (bool, error)
}

// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-computersystem
type HyperVVirtualMachine struct {
	instancePath string
	Uuid         string   `json:"uuid"`
	Name         string   `json:"name"`
	Description  *string  `json:"description"`
	Status       VMStatus `json:"status"`
	CpuCoreCount int      `json:"cpu_core_count"`
	MemorySize   int      `json:"memory_size"`
	SavePath     string   `json:"save_path"`
	*virtualsystem.VirtualMachine
}

func vmStatus(vm *virtualsystem.VirtualMachine) (VMStatus, error) {
	state, err := vm.State()
	if err != nil {
		return VMStatusUnknown, err
	}
	// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/requeststatechange-msvm-computersystem
	if state == virtualsystem.Unknown {
		return VMStatusUnknown, nil
	} else if state == virtualsystem.Starting {
		return VMStatusStarting, nil
	} else if state == virtualsystem.Running {
		return VMStatusRunning, nil
	} else if state == virtualsystem.Off {
		return VMStatusStopped, nil
	} else if state == virtualsystem.Stopping {
		return VMStatusStopping, nil
	} else if state == virtualsystem.Saving {
		return VMStatusSuspending, nil
	} else if state == virtualsystem.Saved {
		return VMStatusSuspended, nil
	}
	return VMStatusUnknown, nil
}

func (vm *HyperVVirtualMachine) update(virtualMachine *virtualsystem.VirtualMachine) (err error) {
	var (
		systemSettingData    *virtualsystem.VirtualSystemSettingData
		processorSettingData *processor.ProcessorSettingData
		cpuCoreCount         uint64
		memorySettingData    *memory.MemorySettingData
		memorySize           uint64
	)
	// 1. 更新虚拟机状态
	if vm.Status, err = vmStatus(virtualMachine); err != nil {
		return
	}
	// 2. 更新虚拟机配置信息
	if systemSettingData, err = virtualMachine.GetVirtualSystemSettingData(); err != nil {
		return
	}
	// 2.1 更新虚拟机保存路径
	if vm.SavePath, err = systemSettingData.GetPropertyConfigurationDataRoot(); err != nil {
		return
	}
	// 2.2 更新虚拟机描述信息
	notes, err := systemSettingData.GetPropertyNotes()
	if err != nil {
		return
	}
	description := slice.Join(notes, ",")
	vm.Description = &description
	// 3. 更新虚拟机 CPU 核心数
	if processorSettingData, err = virtualMachine.GetProcessor(); err != nil {
		return
	}
	if cpuCoreCount, err = processorSettingData.GetCPUCount(); err != nil {
		return
	}
	vm.CpuCoreCount = int(cpuCoreCount)
	// 4. 更新虚拟机内存大小
	if memorySettingData, err = virtualMachine.GetMemory(); err != nil {
		return
	}
	if memorySize, err = memorySettingData.GetSizeMB(); err != nil {
		return
	}
	vm.MemorySize = int(memorySize)
	// 5. 更新虚拟机实例
	vm.VirtualMachine = virtualMachine
	vm.Uuid = virtualMachine.ID()
	vm.Name = virtualMachine.Name()
	vm.instancePath = virtualMachine.InstancePath()
	return
}

func hyperVVirtualMachine(virtualMachine *virtualsystem.VirtualMachine) (*HyperVVirtualMachine, error) {
	var vm HyperVVirtualMachine
	if err := vm.update(virtualMachine); err != nil {
		return nil, err
	}
	return &vm, nil
}

func (vm *HyperVVirtualMachine) VM() (*virtualsystem.VirtualMachine, error) {
	wmiInstance, err := wmiext.GetWmiInstanceFromPath(wmiext.VirtualizationV2, vm.instancePath)
	if err != nil {
		return nil, err
	}
	virtualMachine, err := virtualsystem.NewVirtualMachine(wmiInstance)
	if err != nil {
		return nil, err
	}
	if err = vm.update(virtualMachine); err != nil {
		return nil, err
	}
	return virtualMachine, nil
}

func (vm *HyperVVirtualMachine) Start() (bool, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	// 如果虚拟机正在启动, 则返回启动失败, 并返回错误: 虚拟机正在启动
	if vm.Status == VMStatusStarting {
		return false, errors.New("vm is starting")
	}
	// 如果虚拟机已经启动, 则返回启动成功, 并返回错误: 虚拟机已经启动
	if vm.Status == VMStatusRunning {
		return true, ErrVmAlreadyRunning
	}
	// 如果虚拟机正在关闭, 则返回启动失败, 并返回错误: 虚拟机正在关闭
	if vm.Status == VMStatusStopping {
		return false, errors.New("vm is stopping")
	}
	if err = virtualMachine.Start(); err != nil {
		return false, err
	}
	vm.Status = VMStatusRunning
	return true, nil
}

func (vm *HyperVVirtualMachine) Stop(force bool) (bool, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	// 如果虚拟机已经停止, 则返回关闭成功, 并返回错误: 虚拟机已经停止
	if vm.Status == VMStatusStopped {
		return true, ErrVmAlreadyStopped
	}
	// 如果虚拟机正在启动, 则返回关闭失败, 并返回错误: 虚拟机正在启动
	if !force && vm.Status == VMStatusStarting {
		return false, errors.New("vm is starting")
	}
	// 如果虚拟机正在关闭, 则返回关闭失败, 并返回错误: 虚拟机正在关闭
	if !force && vm.Status == VMStatusStopping {
		return false, errors.New("vm is stopping")
	}
	if err = virtualMachine.Stop(force); err != nil {
		return false, err
	}
	vm.Status = VMStatusStopped
	return true, nil
}

func (vm *HyperVVirtualMachine) Reboot(force bool) (bool, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	// 如果虚拟机正在启动, 则返回重启失败, 并返回错误: 虚拟机正在启动
	if vm.Status == VMStatusStarting {
		return false, errors.New("vm is starting")
	}
	// 如果虚拟机正在关闭, 则返回重启失败, 并返回错误: 虚拟机正在关闭
	if vm.Status == VMStatusStopping {
		return false, errors.New("vm is stopping")
	}
	vm.Status = VMStatusRebooting
	if err = virtualMachine.Stop(force); err != nil {
		return false, err
	}
	vm.Status = VMStatusStopped
	vm.Status = VMStatusStarting
	if err = virtualMachine.Start(); err != nil {
		return false, err
	}
	vm.Status = VMStatusRunning
	return true, nil
}

func (vm *HyperVVirtualMachine) ForceStop() (bool, error) {
	return vm.Stop(true)
}

func (vm *HyperVVirtualMachine) ForceReboot() (bool, error) {
	return vm.Reboot(true)
}

func (vm *HyperVVirtualMachine) Save() (bool, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	// 如果虚拟机已经保存, 则返回保存成功, 并返回错误: 虚拟机已经保存
	//if vm.Status == VMStatusSaved {
	//	return true, ErrVmAlreadySaved
	//}
	if vm.Status == VMStatusSuspended {
		return true, ErrVmAlreadySaved
	}
	// 如果虚拟机不在运行, 则返回保存失败, 并返回错误: 虚拟机不在运行
	if vm.Status != VMStatusRunning {
		return false, errors.New("vm is not running")
	}
	if err = virtualMachine.Save(); err != nil {
		return false, err
	}
	vm.Status = VMStatusSuspended
	return true, nil
}

// Suspend 挂起虚拟机
// 由于 Hyper-V 平台原生挂起功能并非真正意义上的挂起, 而是保存虚拟机的状态, 因此这里的挂起操作实际上是保存虚拟机的状态
// 保存虚拟机的状态后, 可以通过 Resume 恢复虚拟机的运行
func (vm *HyperVVirtualMachine) Suspend() (bool, error) {
	return vm.Save()
}

func (vm *HyperVVirtualMachine) Resume() (bool, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	// 如果虚拟机已经运行, 则返回恢复成功, 并返回错误: 虚拟机已经运行
	if vm.Status == VMStatusRunning {
		return true, ErrVmAlreadyRunning
	}
	// 如果虚拟机不在挂起状态, 则返回恢复失败, 并返回错误: 虚拟机不在挂起状态
	if vm.Status != VMStatusSuspended && vm.Status != VMStatusSaved {
		return false, errors.New("vm is not suspended")
	}

	//if err = virtualMachine.Resume(); err != nil {
	//	return true, err
	//}
	if err = virtualMachine.Start(); err != nil {
		return true, err
	}

	vm.Status = VMStatusRunning
	return true, nil
}

func (vm *HyperVVirtualMachine) Create() (bool, error) {
	var err error

	builder, err := NewVirtualMachineBuilder()
	if err != nil {
		return false, err
	}
	if _, err := builder.SetConfigurationDataRoot(vm.SavePath); err != nil {
		return false, err
	}
	if _, err := builder.SetName(vm.Name); err != nil {
		return false, err
	}
	if _, err := builder.SetCpuCoreCount(uint64(vm.CpuCoreCount)); err != nil {
		return false, err
	}
	if _, err := builder.SetMemorySize(uint64(vm.MemorySize)); err != nil {
		return false, err
	}
	if _, err := builder.SetEnableDynamicMemory(false); err != nil {
		return false, err
	}

	vm.Status = VMStatusCreating
	virtualMachine, err := builder.Create()
	if err != nil {
		if err := errors.Unwrap(err); err != nil && strings.Contains(err.Error(), "ErrorCode[32769]") {
			// 创建失败, 当前目录下已存在同名的虚拟机
			return false, ErrVmAlreadyExists
		}
		return false, err
	}
	if err = vm.update(virtualMachine.VirtualMachine); err != nil {
		return false, err
	}
	vm.Status = VMStatusStopped
	log.Printf("vm created: %v, %v\n", vm, virtualMachine)
	return true, err
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

func (vm *HyperVVirtualMachine) ModifySpec(options ...Option) (ok bool, err error) {
	opts := new(ModifySpecOptions)
	for _, option := range options {
		option(opts)
	}
	virtualMachine, err := vm.VM()
	if err != nil {
		return false, err
	}
	if opts.confirmStop && vm.Status != VMStatusStopped {
		if ok, err = vm.ForceStop(); err != nil {
			return
		}
	}
	mgmt, err := wmiext.NewLocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	// 修改 CPU 核心数
	if opts.cpuCoreCount > 0 {
		// HyperV 要求如果虚拟机未关闭, 则不允许修改 CPU 核心数
		if vm.Status != VMStatusStopped {
			return false, errors.New("vm must be stopped before modifying spec")
		}
		if err = mgmt.SetProcessorCount(virtualMachine, uint64(opts.cpuCoreCount)); err != nil {
			return false, err
		}
	}
	if opts.memorySizeMB > 0 {
		// HyperV 允许虚拟机运行状态下, 修改内存大小
		if err = mgmt.SetMemoryMB(virtualMachine, uint64(opts.memorySizeMB)); err != nil {
			return false, err
		}
	}
	return true, nil
}

// GetVirtualMachineByName 根据虚拟机名称获取虚拟机
func GetVirtualMachineByName(vmName string) (*HyperVVirtualMachine, error) {
	vm, err := wmiext.GetVirtualMachineByVMName(vmName)
	if err != nil {
		return nil, err
	}
	virtualMachine, err := hyperVVirtualMachine(vm)
	if err != nil {
		return nil, err
	}
	return virtualMachine, nil
}

// CreateVirtualMachine 创建虚拟机
func CreateVirtualMachine(name string, savePath string, cpuCoreCount int, memorySize int) (*HyperVVirtualMachine, error) {
	hyperVVirtualMachine := &HyperVVirtualMachine{
		Name:         name,
		Status:       VMStatusCreating,
		SavePath:     savePath,
		CpuCoreCount: cpuCoreCount,
		MemorySize:   memorySize,
	}
	ok, err := hyperVVirtualMachine.Create()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("create virtual machine failed")
	}
	hyperVVirtualMachine.Status = VMStatusStopped
	return hyperVVirtualMachine, nil
}

// ListVirtualMachines 获取所有虚拟机
func ListVirtualMachines() (vms []*HyperVVirtualMachine, err error) {
	vmms, err := wmiext.NewLocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	virtualMachines, err := vmms.GetVirtualMachines()
	if err != nil {
		return nil, err
	}
	for _, virtualMachine := range virtualMachines {
		vm, err := hyperVVirtualMachine(virtualMachine)
		if err != nil {
			return nil, err
		}
		vms = append(vms, vm)
	}
	return
}

// DeleteVirtualMachineByName 根据名称删除虚拟机
func DeleteVirtualMachineByName(name string, del bool) (ok bool, err error) {
	vm, err := GetVirtualMachineByName(name)
	if err != nil {
		return false, err
	}
	virtualMachine := vm.VirtualMachine
	vmms, err := wmiext.NewLocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	//defer vmms.Close()
	if vm.Status != VMStatusStopped {
		return false, errors.New("vm must be stopped before deleting")
	}
	err = vmms.DeleteVirtualMachine(virtualMachine)
	if err != nil {
		return false, err
	}
	if del {
		// RemoveAll 可以删除非空文件夹
		err = os.RemoveAll(vm.SavePath)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

// ModifyVirtualMachineSpec 根据虚拟机名称修改虚拟机规格
func ModifyVirtualMachineSpec(name string, cpuCoreCount int, memorySize int) (ok bool, err error) {
	vm, err := GetVirtualMachineByName(name)
	if err != nil {
		return false, err
	}
	var options []Option
	if cpuCoreCount > 0 {
		options = append(options, WithCpuCoreCount(cpuCoreCount))
		options = append(options, WithStop(true))
	}
	if memorySize > 0 {
		options = append(options, WithMemorySize(memorySize))
	}
	return vm.ModifySpec(options...)
}

func waitVMResult(res int32, service *wmiext.Service, job *wmiext.Instance, errorMsg string, translate func(int) error) error {
	var err error

	switch res {
	case 0:
		return nil
	case 4096:
		err = wmiext.WaitJob(service, job)
		defer job.Close()
	default:
		if translate != nil {
			return translate(int(res))
		}

		return fmt.Errorf("%s (result code %d)", errorMsg, res)
	}

	if err != nil {
		desc, _ := job.GetAsString("ErrorDescription")
		desc = strings.Replace(desc, "\n", " ", -1)
		return fmt.Errorf("%s: %w (%s)", errorMsg, err, desc)
	}

	return err
}
