	package networking

import "github.com/rokukoo/hyperv/pkg/wmiext"

// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-syntheticethernetportsettingdata

const (
	Msvm_SyntheticEthernetPortSettingData = "Msvm_SyntheticEthernetPortSettingData"
)

// SyntheticEthernetPortSettingData 代表 Hyper-V 的虚拟以太网端口默认设置
type SyntheticEthernetPortSettingData struct {
	S__PATH string `json:"-"`

	InstanceID               string   `json:"instance_id"`
	Caption                  string   `json:"caption" default:"Virtual Ethernet Port Default Settings"`
	Description              string   `json:"description" default:"Describes the default settings for the virtual Ethernet port resources."`
	ElementName              string   `json:"element_name"`
	ResourceType             uint16   `json:"resource_type" default:"10"`
	OtherResourceType        string   `json:"other_resource_type"`
	ResourceSubType          string   `json:"resource_sub_type" default:"Microsoft:Hyper-V:Synthetic Ethernet Port"`
	PoolID                   string   `json:"pool_id"`
	ConsumerVisibility       uint16   `json:"consumer_visibility" default:"3"`
	HostResource             []string `json:"host_resource"`
	AllocationUnits          string   `json:"allocation_units" default:"count"`
	VirtualQuantity          uint64   `json:"virtual_quantity" default:"1"`
	Reservation              uint64   `json:"reservation" default:"1"`
	Limit                    uint64   `json:"limit" default:"1"`
	Weight                   uint32   `json:"weight" default:"0"`
	AutomaticAllocation      bool     `json:"automatic_allocation" default:"true"`
	AutomaticDeallocation    bool     `json:"automatic_deallocation" default:"true"`
	Parent                   string   `json:"parent"`
	Connection               []string `json:"connection"`
	Address                  string   `json:"address"`
	MappingBehavior          uint16   `json:"mapping_behavior"`
	AddressOnParent          string   `json:"address_on_parent"`
	VirtualQuantityUnits     string   `json:"virtual_quantity_units" default:"count"`
	DesiredVLANEndpointMode  uint16   `json:"desired_vlan_endpoint_mode"`
	OtherEndpointMode        string   `json:"other_endpoint_mode"`
	VirtualSystemIdentifiers []string `json:"virtual_system_identifiers"`
	DeviceNamingEnabled      bool     `json:"device_naming_enabled" default:"false"`
	AllowPacketDirect        bool     `json:"allow_packet_direct" default:"false"`
	StaticMacAddress         bool     `json:"static_mac_address" default:"false"`
	ClusterMonitored         bool     `json:"cluster_monitored" default:"true"`

	*wmiext.Instance `json:"-"`
}

func (sepsd *SyntheticEthernetPortSettingData) Path() string {
	return sepsd.S__PATH
}
