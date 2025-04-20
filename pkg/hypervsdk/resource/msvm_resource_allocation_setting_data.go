package resource

import (
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	Msvm_ResourceAllocationSettingData = "Msvm_ResourceAllocationSettingData"
)

type SettingsDefineCapabilities_ValueRole int

const (
	// Default enum
	SettingsDefineCapabilities_ValueRole_Default SettingsDefineCapabilities_ValueRole = 0
	// Optimal enum
	SettingsDefineCapabilities_ValueRole_Optimal SettingsDefineCapabilities_ValueRole = 1
	// Mean enum
	SettingsDefineCapabilities_ValueRole_Mean SettingsDefineCapabilities_ValueRole = 2
	// Supported enum
	SettingsDefineCapabilities_ValueRole_Supported SettingsDefineCapabilities_ValueRole = 3
	// DMTF_Reserved enum
	SettingsDefineCapabilities_ValueRole_DMTF_Reserved SettingsDefineCapabilities_ValueRole = 4
)

type SettingsDefineCapabilities_ValueRange int

const (
	// Point enum
	SettingsDefineCapabilities_ValueRange_Point SettingsDefineCapabilities_ValueRange = 0
	// Minimums enum
	SettingsDefineCapabilities_ValueRange_Minimums SettingsDefineCapabilities_ValueRange = 1
	// Maximums enum
	SettingsDefineCapabilities_ValueRange_Maximums SettingsDefineCapabilities_ValueRange = 2
	// Increments enum
	SettingsDefineCapabilities_ValueRange_Increments SettingsDefineCapabilities_ValueRange = 3
	// DMTF_Reserved enum
	SettingsDefineCapabilities_ValueRange_DMTF_Reserved SettingsDefineCapabilities_ValueRange = 4
)

type ResourceAllocationSettingData struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID               string   `json:"instance_id"`
	Caption                  string   `json:"caption"`
	Description              string   `json:"description"`
	ElementName              string   `json:"element_name"`
	ResourceType             uint16   `json:"resource_type"`
	OtherResourceType        string   `json:"other_resource_type"`
	ResourceSubType          string   `json:"resource_sub_type"`
	PoolID                   string   `json:"pool_id"`
	ConsumerVisibility       uint16   `json:"consumer_visibility"`
	HostResource             []string `json:"host_resource"`
	AllocationUnits          string   `json:"allocation_units"`
	VirtualQuantity          uint64   `json:"virtual_quantity"`
	Reservation              uint64   `json:"reservation"`
	Limit                    uint64   `json:"limit"`
	Weight                   uint32   `json:"weight"`
	AutomaticAllocation      bool     `json:"automatic_allocation"`
	AutomaticDeallocation    bool     `json:"automatic_deallocation"`
	Parent                   string   `json:"parent"`
	Connection               []string `json:"connection"`
	Address                  string   `json:"address"`
	MappingBehavior          uint16   `json:"mapping_behavior"`
	AddressOnParent          string   `json:"address_on_parent"`
	VirtualQuantityUnits     string   `json:"virtual_quantity_units"`
	VirtualSystemIdentifiers []string `json:"virtual_system_identifiers"`

	*wmiext.Instance `json:"-"`
}

func (rasd *ResourceAllocationSettingData) Path() string {
	return rasd.S__PATH
}

func (rasd *ResourceAllocationSettingData) SetParent(parent string) error {
	rasd.Parent = parent
	return rasd.Put("Parent", parent)
}

func (rasd *ResourceAllocationSettingData) GetParent() string {
	return rasd.Parent
}

func (rasd *ResourceAllocationSettingData) GetParenObject() (*ResourceAllocationSettingData, error) {
	parentPath := rasd.GetParent()
	if parentPath == "" {
		return nil, wmiext.NotFound
	}
	resourceAllocationSettingData := &ResourceAllocationSettingData{}
	return resourceAllocationSettingData, rasd.GetService().GetObjectAsObject(parentPath, resourceAllocationSettingData)
}

func (rasd *ResourceAllocationSettingData) SetAddressOnParent(addressOnParent string) error {
	rasd.AddressOnParent = addressOnParent
	return rasd.Put("AddressOnParent", addressOnParent)
}

func (rasd *ResourceAllocationSettingData) GetResourceType() (rtype *ResourceTypeValue, err error) {
	rsub := rasd.ResourceSubType
	rothersub := rasd.OtherResourceType

	resourceType := rasd.ResourceType
	// v2.ResourceAllocationSettingData_ResourceType

	rtype = &ResourceTypeValue{
		ResourceType:      ResourcePool_ResourceType(int(resourceType)),
		OtherResourceType: rothersub,
		ResourceSubType:   rsub,
	}
	return
}
