package hypervctl

import (
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/memory"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/processor"
	virtualsystem "github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system"
)

type VirtualMachineBuilder struct {
	// name                string
	// cpuCoreCount        int
	// memorySize          int
	// enableDynamicMemory bool

	Err error

	systemSettingsData *virtualsystem.VirtualSystemSettingData
	memorySettingData  *memory.MemorySettingsData
	processorSettings  *processor.ProcessorSettingData

	svc *virtualsystem.VirtualSystemManagementService
}

func NewVirtualMachineBuilder() (*VirtualMachineBuilder, error) {
	return &VirtualMachineBuilder{}, nil
}

func (builder *VirtualMachineBuilder) PrepareSystemSettings(name string, beforeAdd func(systemSettingsData *virtualsystem.VirtualSystemSettingData)) *VirtualMachineBuilder {
	if builder.Err != nil {
		return builder
	}

	if builder.systemSettingsData == nil {
		builder.systemSettingsData = virtualsystem.GetDefaultVirtualSystemSettingData()
		builder.systemSettingsData.ElementName = name
		builder.systemSettingsData.VirtualSystemSubType = virtualsystem.HyperVGeneration_V1
		builder.systemSettingsData.AutomaticSnapshotsEnabled = false
	}

	if beforeAdd != nil {
		beforeAdd(builder.systemSettingsData)
	}

	return builder
}

func (builder *VirtualMachineBuilder) PrepareProcessorSettings(beforeAdd func(processorSettings *processor.ProcessorSettingData)) *VirtualMachineBuilder {
	if builder.Err != nil {
		return builder
	}

	if builder.processorSettings == nil {
		if builder.processorSettings, builder.Err = processor.GetDefaultProcessorSettingData(); builder.Err != nil {
			return builder
		}
	}

	if beforeAdd != nil {
		beforeAdd(builder.processorSettings)
	}

	return builder
}

func (builder *VirtualMachineBuilder) PrepareMemorySettings(beforeAdd func(data *memory.MemorySettingsData)) *VirtualMachineBuilder {
	if builder.Err != nil {
		return builder
	}

	if builder.memorySettingData == nil {
		if builder.memorySettingData, builder.Err = memory.GetDefaultMemorySettingsData(); builder.Err != nil {
			return builder
		}
	}

	if beforeAdd != nil {
		beforeAdd(builder.memorySettingData)
	}

	return builder
}

//func NewVirtualMachineBuilder() (*VirtualMachineBuilder, error) {
//	var err error
//	var conn *host.WmiHost
//	var settingsData *virtualsystem2.VirtualSystemSettingData
//	var memorySettingsData *memory.MemorySettingData
//	var processorSettings *processor.ProcessorSettingData
//
//	conn = host.NewWmiLocalHost()
//
//	// Virtual system settings
//	if settingsData, err = virtualsystem2.GetDefaultVirtualSystemSettingData(conn); err != nil {
//		return nil, err
//	}
//	if err = settingsData.SetHyperVGeneration(virtualsystem2.HyperVGeneration_V1); err != nil {
//		return nil, err
//	}
//
//	// Memory settings
//	if memorySettingsData, err = memory.GetDefaultMemorySettingData(conn); err != nil {
//		return nil, err
//	}
//	if err = memorySettingsData.SetPropertyDynamicMemoryEnabled(false); err != nil {
//		return nil, err
//	}
//
//	// Processor settings
//	if processorSettings, err = processor.GetDefaultProcessorSettingData(conn); err != nil {
//		return nil, err
//	}
//
//	return &VirtualMachineBuilder{
//		conn:              conn,
//		settingsData:      settingsData,
//		memorySettingData: memorySettingsData,
//		processorSettings: processorSettings,
//	}, err
//}

//	func (builder *VirtualMachineBuilder) SetGenerationVersion(version virtualsystem.HyperVGeneration) (*VirtualMachineBuilder, error) {
//		if builder.Err != nil {
//			return builder, builder.Err
//		}
//
//		err := builder.settingsData.SetHyperVGeneration(version)
//		builder.setErr(err)
//
//		return builder, nil
//	}
//
//	func (builder *VirtualMachineBuilder) Close() error {
//		var err error
//
//		if builder.settingsData != nil {
//			err = builder.settingsData.Close()
//		}
//
//		if builder.memorySettingData != nil {
//			err = builder.memorySettingData.Close()
//		}
//
//		if builder.processorSettings != nil {
//			err = builder.processorSettings.Close()
//		}
//
//		return err
//	}

func (builder *VirtualMachineBuilder) Build() (*VirtualMachine, error) {
	var err error
	var vmms *virtualsystem.VirtualSystemManagementService
	var cs *virtualsystem.ComputerSystem

	if vmms, err = virtualsystem.LocalVirtualSystemManagementService(); err != nil {
		return nil, err
	}
	builder.svc = vmms
	// defer vmms.Close()
	cs, err = vmms.DefineSystem(builder.systemSettingsData, builder.processorSettings, builder.memorySettingData)

	return NewVirtualMachine(cs)
}

//func (builder *VirtualMachineBuilder) Create() (*virtualsystem.ComputerSystem, error) {
//	var vm *VirtualMachine
//	var err error
//
//	vmms, err := wmiext.NewLocalVirtualSystemManagementService()
//	if err != nil {
//		return nil, err
//	}
//	// defer vmms.Close()
//	machine, err := vmms.CreateVirtualMachine(builder.settingsData, builder.memorySettingData, builder.processorSettings)
//	if err != nil {
//		return nil, err
//	}
//	vm, err = hyperVVirtualMachine(machine)
//	if err != nil {
//		return nil, err
//	}
//	return vm, err
//}

//
//	func (builder *VirtualMachineBuilder) SetEnableDynamicMemory(enable bool) (*VirtualMachineBuilder, error) {
//		if builder.Err != nil {
//			return builder, builder.Err
//		}
//
//		builder.enableDynamicMemory = enable
//		err := builder.memorySettingData.SetPropertyDynamicMemoryEnabled(enable)
//		builder.setErr(err)
//
//		return builder, nil
//	}
//
//	func (builder *VirtualMachineBuilder) SetCpuCoreCount(count uint64) (*VirtualMachineBuilder, error) {
//		if builder.Err != nil {
//			return builder, builder.Err
//		}
//
//		builder.cpuCoreCount = count
//
//		err := builder.processorSettings.SetCPUCount(count)
//		builder.setErr(err)
//
//		return builder, nil
//	}
//func (builder *VirtualMachineBuilder) SetMemorySize(sizeMB uint64) (*VirtualMachineBuilder, error) {
//	if builder.Err != nil {
//		return builder, builder.Err
//	}
//
//	builder.memorySize = sizeMB
//
//	err := builder.memorySettingData.SetSizeMB(sizeMB)
//	builder.setErr(err)
//
//	return builder, nil
//}

//
//func (builder *VirtualMachineBuilder) SetConfigurationDataRoot(path string) (*VirtualMachineBuilder, error) {
//	if builder.Err != nil {
//		return builder, builder.Err
//	}
//
//	err := builder.settingsData.SetPropertyConfigurationDataRoot(path)
//	builder.setErr(err)
//
//	return builder, nil
//}
//
//func (builder *VirtualMachineBuilder) SetName(s string) (*VirtualMachineBuilder, error) {
//	if builder.Err != nil {
//		return builder, builder.Err
//	}
//
//	builder.name = s
//
//	err := builder.settingsData.SetPropertyElementName(s)
//	builder.setErr(err)
//
//	return builder, nil
//}
//
//func (builder *VirtualMachineBuilder) setErr(err error) {
//	builder.Err = err
//}
