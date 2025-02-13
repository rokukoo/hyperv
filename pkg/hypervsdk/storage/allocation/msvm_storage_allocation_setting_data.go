package allocation

import "github.com/rokukoo/hypervctl/pkg/wmiext"

const (
	Msvm_StorageAllocationSettingData = "Msvm_StorageAllocationSettingData "
)

type StorageAllocationSettingData struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID                      string   `json:"instance_id"`
	Caption                         string   `json:"caption"`
	Description                     string   `json:"description"`
	ElementName                     string   `json:"element_name"`
	ResourceType                    uint16   `json:"resource_type"`
	OtherResourceType               string   `json:"other_resource_type"`
	ResourceSubType                 string   `json:"resource_sub_type"`
	PoolID                          string   `json:"pool_id"`
	ConsumerVisibility              uint16   `json:"consumer_visibility"`
	HostResource                    []string `json:"host_resource"`
	AllocationUnits                 string   `json:"allocation_units"`
	VirtualQuantity                 uint64   `json:"virtual_quantity"`
	Limit                           uint64   `json:"limit"`
	Weight                          uint32   `json:"weight"`
	StorageQoSPolicyID              string   `json:"storage_qos_policy_id"`
	AutomaticAllocation             bool     `json:"automatic_allocation"`
	AutomaticDeallocation           bool     `json:"automatic_deallocation"`
	Parent                          string   `json:"parent"`
	Connection                      []string `json:"connection"`
	Address                         string   `json:"address"`
	MappingBehavior                 uint16   `json:"mapping_behavior"`
	AddressOnParent                 string   `json:"address_on_parent"`
	VirtualResourceBlockSize        uint64   `json:"virtual_resource_block_size"`
	VirtualQuantityUnits            string   `json:"virtual_quantity_units"`
	Access                          uint16   `json:"access"`
	HostResourceBlockSize           uint64   `json:"host_resource_block_size"`
	Reservation                     uint64   `json:"reservation"`
	HostExtentStartingAddress       uint64   `json:"host_extent_starting_address"`
	HostExtentName                  string   `json:"host_extent_name"`
	HostExtentNameFormat            uint16   `json:"host_extent_name_format"`
	OtherHostExtentNameFormat       string   `json:"other_host_extent_name_format"`
	HostExtentNameNamespace         uint16   `json:"host_extent_name_namespace"`
	OtherHostExtentNameNamespace    string   `json:"other_host_extent_name_namespace"`
	IOPSLimit                       uint64   `json:"iops_limit"`
	IOPSReservation                 uint64   `json:"iops_reservation"`
	IOPSAllocationUnits             string   `json:"iops_allocation_units"`
	PersistentReservationsSupported bool     `json:"persistent_reservations_supported"`
	CachingMode                     uint16   `json:"caching_mode"`
	SnapshotId                      string   `json:"snapshot_id"`
	IgnoreFlushes                   bool     `json:"ignore_flushes"`
	WriteHardeningMethod            uint16   `json:"write_hardening_method"`

	*wmiext.Instance `json:"-"`
}

func (s *StorageAllocationSettingData) Path() string {
	return s.S__PATH
}

func (s *StorageAllocationSettingData) SetParent(parent string) (err error) {
	s.Parent = parent
	return s.Put("Parent", parent)
}

func (s *StorageAllocationSettingData) SetHostResource(hostResource []string) (err error) {
	s.HostResource = hostResource
	return s.Put("HostResource", hostResource)
}
