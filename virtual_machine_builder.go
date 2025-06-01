package hyperv

import (
	"github.com/rokukoo/hyperv/pkg/hypervsdk/memory"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/processor"
	virtualsystem "github.com/rokukoo/hyperv/pkg/hypervsdk/virtual_system"
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
	if err != nil {
		return nil, err
	}

	// 提前创建SCSI控制器, 以支持热插拔硬盘
	if err = vmms.AddSCSIController(cs); err != nil {
		return nil, err
	}
	return NewVirtualMachine(cs)
}
