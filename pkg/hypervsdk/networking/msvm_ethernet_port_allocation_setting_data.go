package networking

import (
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking/switch_extension"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	Msvm_EthernetPortAllocationSettingData = "Msvm_EthernetPortAllocationSettingData"
)

type EthernetPortAllocationSettingData struct {
	S__PATH                 string `json:"-"`
	S__CLASS                string `json:"-"`
	InstanceID              string
	Caption                 string
	Description             string
	ElementName             string
	ResourceType            uint16
	OtherResourceType       string
	ResourceSubType         string
	PoolID                  string
	ConsumerVisibility      uint16
	HostResource            []string
	AllocationUnits         string
	VirtualQuantity         uint64
	Reservation             uint64
	Limit                   uint64
	Weight                  uint32
	AutomaticAllocation     bool
	AutomaticDeallocation   bool
	Parent                  string
	Connection              []string
	Address                 string
	MappingBehavior         uint16
	AddressOnParent         string
	VirtualQuantityUnits    string
	DesiredVLANEndpointMode uint16
	OtherEndpointMode       string
	EnabledState            uint16
	LastKnownSwitchName     string
	RequiredFeatures        []string
	RequiredFeatureHints    []string
	TestReplicaPoolID       string
	TestReplicaSwitchName   string
	CompartmentGuid         string

	*wmiext.Instance
}

func (epasd *EthernetPortAllocationSettingData) Path() string {
	return epasd.S__PATH
}

func NewEthernetPortAllocationSettingDataFromInstance(instance *wmiext.Instance) (*EthernetPortAllocationSettingData, error) {
	epasd := &EthernetPortAllocationSettingData{}
	if err := instance.GetAll(epasd); err != nil {
		return nil, err
	}
	return epasd, nil
}

func (epasd *EthernetPortAllocationSettingData) SetEnabledState(enabledState uint16) error {
	epasd.EnabledState = enabledState
	return epasd.Put("EnabledState", enabledState)
}

func (epasd *EthernetPortAllocationSettingData) SetHostResource(hostResource []string) error {
	epasd.HostResource = hostResource
	return epasd.Put("HostResource", hostResource)
}

func (epasd *EthernetPortAllocationSettingData) GetEthernetSwitchPortBandwidthSettingData() (*switch_extension.EthernetSwitchPortBandwidthSettingData, error) {
	inst, err := epasd.GetRelated(switch_extension.Msvm_EthernetSwitchPortBandwidthSettingData)
	if err != nil {
		return nil, err
	}
	return switch_extension.NewEthernetSwitchPortBandwidthSettingData(inst)
}

func (epasd *EthernetPortAllocationSettingData) GetEthernetSwitchPortVlanSettingData() (*switch_extension.EthernetSwitchPortVlanSettingData, error) {
	inst, err := epasd.GetRelated(switch_extension.Msvm_EthernetSwitchPortVlanSettingData)
	if err != nil {
		return nil, err
	}
	return switch_extension.NewEthernetSwitchPortVlanSettingData(inst)
}
