package hypervctl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/memory"
	"github.com/microsoft/wmi/pkg/virtualization/core/processor"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	"github.com/rokukoo/hypervctl/wmiext"
)

type VirtualMachineBuilder struct {
	name                string
	cpuCoreCount        uint64
	memorySize          uint64
	enableDynamicMemory bool

	Err               error
	conn              *host.WmiHost
	settingsData      *virtualsystem.VirtualSystemSettingData
	memorySettingData *memory.MemorySettingData
	processorSettings *processor.ProcessorSettingData
}

func NewVirtualMachineBuilder() (*VirtualMachineBuilder, error) {
	var err error
	var conn *host.WmiHost
	var settingsData *virtualsystem.VirtualSystemSettingData
	var memorySettingsData *memory.MemorySettingData
	var processorSettings *processor.ProcessorSettingData

	conn = host.NewWmiLocalHost()

	// Virtual system settings
	if settingsData, err = virtualsystem.GetDefaultVirtualSystemSettingData(conn); err != nil {
		return nil, err
	}
	if err = settingsData.SetHyperVGeneration(virtualsystem.HyperVGeneration_V1); err != nil {
		return nil, err
	}

	// Memory settings
	if memorySettingsData, err = memory.GetDefaultMemorySettingData(conn); err != nil {
		return nil, err
	}
	if err = memorySettingsData.SetPropertyDynamicMemoryEnabled(false); err != nil {
		return nil, err
	}

	// Processor settings
	if processorSettings, err = processor.GetDefaultProcessorSettingData(conn); err != nil {
		return nil, err
	}

	return &VirtualMachineBuilder{
		conn:              conn,
		settingsData:      settingsData,
		memorySettingData: memorySettingsData,
		processorSettings: processorSettings,
	}, err
}

func (builder *VirtualMachineBuilder) SetGenerationVersion(version virtualsystem.HyperVGeneration) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	err := builder.settingsData.SetHyperVGeneration(version)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) Close() error {
	var err error

	if builder.settingsData != nil {
		err = builder.settingsData.Close()
	}

	if builder.memorySettingData != nil {
		err = builder.memorySettingData.Close()
	}

	if builder.processorSettings != nil {
		err = builder.processorSettings.Close()
	}

	return err
}

func (builder *VirtualMachineBuilder) Create() (*HyperVVirtualMachine, error) {
	var vm *HyperVVirtualMachine
	var err error

	vmms, err := wmiext.NewLocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	// defer vmms.Close()
	machine, err := vmms.CreateVirtualMachine(builder.settingsData, builder.memorySettingData, builder.processorSettings)
	if err != nil {
		return nil, err
	}
	vm, err = hyperVVirtualMachine(machine)
	if err != nil {
		return nil, err
	}
	return vm, err
}

func (builder *VirtualMachineBuilder) SetEnableDynamicMemory(enable bool) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	builder.enableDynamicMemory = enable
	err := builder.memorySettingData.SetPropertyDynamicMemoryEnabled(enable)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) SetCpuCoreCount(count uint64) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	builder.cpuCoreCount = count

	err := builder.processorSettings.SetCPUCount(count)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) SetMemorySize(sizeMB uint64) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	builder.memorySize = sizeMB

	err := builder.memorySettingData.SetSizeMB(sizeMB)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) SetConfigurationDataRoot(path string) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	err := builder.settingsData.SetPropertyConfigurationDataRoot(path)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) SetName(s string) (*VirtualMachineBuilder, error) {
	if builder.Err != nil {
		return builder, builder.Err
	}

	builder.name = s

	err := builder.settingsData.SetPropertyElementName(s)
	builder.setErr(err)

	return builder, nil
}

func (builder *VirtualMachineBuilder) setErr(err error) {
	builder.Err = err
}
