package virtual_system

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/memory"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/network_adapter"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/processor"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/allocation"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

const (
	Msvm_VirtualSystemSettingData = "Msvm_VirtualSystemSettingData"
)

type VirtualSystemSettingData struct {
	S__PATH                              string
	InstanceID                           string
	Caption                              string // = "Virtual Machine Settings"
	Description                          string
	ElementName                          string
	VirtualSystemIdentifier              string
	VirtualSystemType                    string
	Notes                                []string
	CreationTime                         time.Time
	ConfigurationID                      string
	ConfigurationDataRoot                string
	ConfigurationFile                    string
	SnapshotDataRoot                     string
	SuspendDataRoot                      string
	SwapFileDataRoot                     string
	LogDataRoot                          string
	AutomaticStartupAction               uint16 // non-zero
	AutomaticStartupActionDelay          time.Duration
	AutomaticStartupActionSequenceNumber uint16
	AutomaticShutdownAction              uint16 // non-zero
	AutomaticRecoveryAction              uint16 // non-zero
	RecoveryFile                         string
	BIOSGUID                             string
	BIOSSerialNumber                     string
	BaseBoardSerialNumber                string
	ChassisSerialNumber                  string
	Architecture                         string
	ChassisAssetTag                      string
	BIOSNumLock                          bool
	BootOrder                            []uint16
	Parent                               string
	UserSnapshotType                     uint16 // non-zero
	IsSaved                              bool
	AdditionalRecoveryInformation        string
	AllowFullSCSICommandSet              bool
	DebugChannelId                       uint32
	DebugPortEnabled                     uint16
	DebugPort                            uint32
	Version                              string
	IncrementalBackupEnabled             bool
	VirtualNumaEnabled                   bool
	AllowReducedFcRedundancy             bool // = False
	VirtualSystemSubType                 string
	BootSourceOrder                      []string
	PauseAfterBootFailure                bool
	NetworkBootPreferredProtocol         uint16 // non-zero
	GuestControlledCacheTypes            bool
	AutomaticSnapshotsEnabled            bool
	IsAutomaticSnapshot                  bool
	GuestStateFile                       string
	GuestStateDataRoot                   string
	LockOnDisconnect                     bool
	ParentPackage                        string
	AutomaticCriticalErrorActionTimeout  time.Duration
	AutomaticCriticalErrorAction         uint16
	ConsoleMode                          uint16
	SecureBootEnabled                    bool
	SecureBootTemplateId                 string
	LowMmioGapSize                       uint64
	HighMmioGapSize                      uint64
	EnhancedSessionTransportType         uint16

	*wmiext.Instance `json:"-"`
}

func (vssd *VirtualSystemSettingData) Path() string {
	return vssd.S__PATH
}

func (vssd *VirtualSystemSettingData) GetStorageAllocationSettingData() (col []*allocation.StorageAllocationSettingData, err error) {
	var (
		resourceAllocationSettingData *allocation.StorageAllocationSettingData
	)
	resourceAllocationSettingDatas, err := vssd.GetAllRelated("Msvm_StorageAllocationSettingData")
	if err != nil {
		return nil, err
	}

	for _, resourceAllocationSettingDataInst := range resourceAllocationSettingDatas {
		resourceAllocationSettingData = &allocation.StorageAllocationSettingData{}
		if err = resourceAllocationSettingDataInst.GetAll(resourceAllocationSettingData); err != nil {
			return
		}
		col = append(col, resourceAllocationSettingData)
	}
	return
}

func (vssd *VirtualSystemSettingData) getResourceAllocationSettingData(rtype resource.ResourceAllocationSettingData_ResourceType) (col []*resource.ResourceAllocationSettingData, err error) {
	var (
		resourceAllocationSettingData *resource.ResourceAllocationSettingData
	)
	resourceType := uint16(rtype)
	resourceAllocationSettingDatas, err := vssd.GetAllRelated("Msvm_ResourceAllocationSettingData")
	for _, resourceAllocationSettingDataInst := range resourceAllocationSettingDatas {
		resourceAllocationSettingData = &resource.ResourceAllocationSettingData{}
		if err = resourceAllocationSettingDataInst.GetAll(resourceAllocationSettingData); err != nil {
			return
		}

		if resourceAllocationSettingData.ResourceType == resourceType {
			col = append(col, resourceAllocationSettingData)
		}
	}
	return
}

// GetComputerSystem returns the ComputerSystem instance that this VirtualSystemSettingData instance is associated with.
func (vssd *VirtualSystemSettingData) GetComputerSystem() (*ComputerSystem, error) {
	var system = &ComputerSystem{}
	var err error

	if err = vssd.GetService().FindFirstRelatedObject(vssd.Path(), "Msvm_ComputerSystem", system); err != nil {
		return nil, err
	}

	return system, nil
}

func (vssd *VirtualSystemSettingData) GetProcessorSettingData() (*processor.ProcessorSettingData, error) {
	var processorSettingData = processor.ProcessorSettingData{}
	return &processorSettingData, vssd.GetService().FindFirstRelatedObject(vssd.Path(), processor.Msvm_ProcessorSettingData, &processorSettingData)
}

func (vssd *VirtualSystemSettingData) GetMemorySettingsData() (*memory.MemorySettingsData, error) {
	var memorySettingsData = memory.MemorySettingsData{}
	return &memorySettingsData, vssd.GetService().FindFirstRelatedObject(vssd.Path(), memory.Msvm_MemorySettingData, &memorySettingsData)
}

func GetDefaultVirtualSystemSettingData() *VirtualSystemSettingData {
	return &VirtualSystemSettingData{}
}

// TODO: Implement this method
func (vsms *VirtualSystemManagementService) CreateSystemSettings(settings *VirtualSystemSettingData) (string, error) {
	systemSettingsInst, err := vsms.Session.SpawnInstance(Msvm_VirtualSystemSettingData)
	if err != nil {
		return "", err
	}

	if err = systemSettingsInst.Put("ElementName", settings.ElementName); err != nil {
		return "", err
	}

	if err = systemSettingsInst.Put("VirtualSystemSubType", settings.VirtualSystemSubType); err != nil {
		return "", err
	}

	if settings.ConfigurationDataRoot != "" {
		if err = systemSettingsInst.Put("ConfigurationDataRoot", settings.ConfigurationDataRoot); err != nil {
			return "", err
		}
	}

	if err = systemSettingsInst.Put("AutomaticSnapshotsEnabled", settings.AutomaticSnapshotsEnabled); err != nil {
		return "", err
	}

	if settings.Notes != nil {
		if err = systemSettingsInst.Put("Notes", settings.Notes); err != nil {
			return "", err
		}
	}

	return systemSettingsInst.GetCimText(), nil
}

func (vsms *VirtualSystemManagementService) AddVirtualEthernetConnection(
	computerSystem *ComputerSystem,
	networkAdapter *network_adapter.VirtualNetworkAdapter,
) (
	epas *networking.EthernetPortAllocationSettingData,
	err error,
) {
	ethernetPortAllocationSettingData, err := computerSystem.NewEthernetPortAllocationSettingData(networkAdapter)
	if err != nil {
		return
	}
	defer ethernetPortAllocationSettingData.Close()

	virtualSystemSettingData, err := computerSystem.GetVirtualSystemSettingData()
	if err != nil {
		return
	}
	defer virtualSystemSettingData.Close()

	// apply the settings
	resultInstances, err := vsms.AddResourceSettings(virtualSystemSettingData, []string{ethernetPortAllocationSettingData.GetCimText()})
	if err != nil {
		return
	}

	if len(resultInstances) == 0 {
		err = errors.Wrapf(wmiext.NotFound, "AddVirtualSystemResource")
		return
	}

	resultInstance, err := resultInstances[0].CloneInstance()
	if err != nil {
		return
	}
	return networking.NewEthernetPortAllocationSettingDataFromInstance(resultInstance)
}

// GetVirtualNetworkAdapterByName
func (vssd *VirtualSystemSettingData) GetVirtualNetworkAdapterByName(adapterName string) (vna *network_adapter.VirtualNetworkAdapter, err error) {
	nas, err := vssd.GetVirtualNetworkAdapters()
	if err != nil {
		return
	}

	for _, networkAdapter := range nas {
		var name string
		name = networkAdapter.ElementName
		if name == adapterName {
			// Found the match
			// Assumption is only one adapter would match, - FIXME - Duplicates
			return networkAdapter.Clone()
		}
	}
	err = errors.Wrapf(wmiext.NotFound, "Virtual Network Adapter with name [%s]", adapterName)
	return
}

func (vssd *VirtualSystemSettingData) GetVirtualNetworkAdapters() (col []*network_adapter.VirtualNetworkAdapter, err error) {
	psds, err := vssd.GetSyntheticVirtualNetworkAdapters()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return nil, err
	}
	psdse, err := vssd.GetEmulatedVirtualNetworkAdapters()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return nil, err
	}

	col = append(col, psds...)
	col = append(col, psdse...)
	return
}

func (vssd *VirtualSystemSettingData) GetSyntheticVirtualNetworkAdapters() (virtualNetworkAdapters []*network_adapter.VirtualNetworkAdapter, err error) {
	syntheticEthernetPortSettingDatas, err := vssd.GetAllRelated("Msvm_SyntheticEthernetPortSettingData")
	if err != nil {
		return nil, err
	}

	virtualNetworkAdapters = make([]*network_adapter.VirtualNetworkAdapter, len(syntheticEthernetPortSettingDatas))
	for i, psd := range syntheticEthernetPortSettingDatas {
		virtualNetworkAdapters[i], err = network_adapter.NewVirtualNetworkAdapterFromInstance(psd)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (vssd *VirtualSystemSettingData) GetEmulatedVirtualNetworkAdapters() (virtualNetworkAdapters []*network_adapter.VirtualNetworkAdapter, err error) {
	syntheticEthernetPortSettingDatas, err := vssd.GetAllRelated("Msvm_EmulatedEthernetPortSettingData")
	if err != nil {
		return nil, err
	}

	virtualNetworkAdapters = make([]*network_adapter.VirtualNetworkAdapter, len(syntheticEthernetPortSettingDatas))
	for i, psd := range syntheticEthernetPortSettingDatas {
		virtualNetworkAdapters[i], err = network_adapter.NewVirtualNetworkAdapterFromInstance(psd)
		if err != nil {
			return nil, err
		}
	}

	return
}
