package resource

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rokukoo/hyperv/pkg/wmiext"
	"time"
)

const (
	Msvm_ResourcePool = "Msvm_ResourcePool"
)

type ResourcePool struct {
	InstanceID                string    `json:"instance_id"`
	Caption                   string    `json:"caption"`
	Description               string    `json:"description"`
	ElementName               string    `json:"element_name"`
	InstallDate               time.Time `json:"install_date"`
	Name                      string    `json:"name"`
	OperationalStatus         []uint16  `json:"operational_status"`
	StatusDescriptions        []string  `json:"status_descriptions"`
	Status                    string    `json:"status"`
	HealthState               uint16    `json:"health_state"`
	CommunicationStatus       uint16    `json:"communication_status"`
	DetailedStatus            uint16    `json:"detailed_status"`
	OperatingStatus           uint16    `json:"operating_status"`
	PrimaryStatus             uint16    `json:"primary_status"`
	PoolID                    string    `json:"pool_id"`
	Primordial                bool      `json:"primordial"`
	Capacity                  uint64    `json:"capacity"`
	Reserved                  uint64    `json:"reserved"`
	ResourceType              uint16    `json:"resource_type"`
	OtherResourceType         string    `json:"other_resource_type"`
	ResourceSubType           string    `json:"resource_sub_type"`
	AllocationUnits           string    `json:"allocation_units"`
	ConsumedResourceUnits     string    `json:"consumed_resource_units"`
	CurrentlyConsumedResource uint64    `json:"currently_consumed_resource"`
	MaxConsumableResource     uint64    `json:"max_consumable_resource"`

	*wmiext.Instance `json:"-"`
}

type ResourcePool_ResourceType int

const (
	// Other enum
	ResourcePool_ResourceType_Other ResourcePool_ResourceType = 1
	// Computer_System enum
	ResourcePool_ResourceType_Computer_System ResourcePool_ResourceType = 2
	// Processor enum
	ResourcePool_ResourceType_Processor ResourcePool_ResourceType = 3
	// Memory enum
	ResourcePool_ResourceType_Memory ResourcePool_ResourceType = 4
	// IDE_Controller enum
	ResourcePool_ResourceType_IDE_Controller ResourcePool_ResourceType = 5
	// Parallel_SCSI_HBA enum
	ResourcePool_ResourceType_Parallel_SCSI_HBA ResourcePool_ResourceType = 6
	// FC_HBA enum
	ResourcePool_ResourceType_FC_HBA ResourcePool_ResourceType = 7
	// iSCSI_HBA enum
	ResourcePool_ResourceType_iSCSI_HBA ResourcePool_ResourceType = 8
	// IB_HCA enum
	ResourcePool_ResourceType_IB_HCA ResourcePool_ResourceType = 9
	// Ethernet_Adapter enum
	ResourcePool_ResourceType_Ethernet_Adapter ResourcePool_ResourceType = 10
	// Other_Network_Adapter enum
	ResourcePool_ResourceType_Other_Network_Adapter ResourcePool_ResourceType = 11
	// I_O_Slot enum
	ResourcePool_ResourceType_I_O_Slot ResourcePool_ResourceType = 12
	// I_O_Device enum
	ResourcePool_ResourceType_I_O_Device ResourcePool_ResourceType = 13
	// Floppy_Drive enum
	ResourcePool_ResourceType_Floppy_Drive ResourcePool_ResourceType = 14
	// CD_Drive enum
	ResourcePool_ResourceType_CD_Drive ResourcePool_ResourceType = 15
	// DVD_drive enum
	ResourcePool_ResourceType_DVD_drive ResourcePool_ResourceType = 16
	// Disk_Drive enum
	ResourcePool_ResourceType_Disk_Drive ResourcePool_ResourceType = 17
	// Tape_Drive enum
	ResourcePool_ResourceType_Tape_Drive ResourcePool_ResourceType = 18
	// Storage_Extent enum
	ResourcePool_ResourceType_Storage_Extent ResourcePool_ResourceType = 19
	// Other_storage_device enum
	ResourcePool_ResourceType_Other_storage_device ResourcePool_ResourceType = 20
	// Serial_port enum
	ResourcePool_ResourceType_Serial_port ResourcePool_ResourceType = 21
	// Parallel_port enum
	ResourcePool_ResourceType_Parallel_port ResourcePool_ResourceType = 22
	// USB_Controller enum
	ResourcePool_ResourceType_USB_Controller ResourcePool_ResourceType = 23
	// Graphics_controller enum
	ResourcePool_ResourceType_Graphics_controller ResourcePool_ResourceType = 24
	// IEEE_1394_Controller enum
	ResourcePool_ResourceType_IEEE_1394_Controller ResourcePool_ResourceType = 25
	// Partitionable_Unit enum
	ResourcePool_ResourceType_Partitionable_Unit ResourcePool_ResourceType = 26
	// Base_Partitionable_Unit enum
	ResourcePool_ResourceType_Base_Partitionable_Unit ResourcePool_ResourceType = 27
	// Power enum
	ResourcePool_ResourceType_Power ResourcePool_ResourceType = 28
	// Cooling_Capacity enum
	ResourcePool_ResourceType_Cooling_Capacity ResourcePool_ResourceType = 29
	// Ethernet_Switch_Port enum
	ResourcePool_ResourceType_Ethernet_Switch_Port ResourcePool_ResourceType = 30
	// Logical_Disk enum
	ResourcePool_ResourceType_Logical_Disk ResourcePool_ResourceType = 31
	// Storage_Volume enum
	ResourcePool_ResourceType_Storage_Volume ResourcePool_ResourceType = 32
	// Ethernet_Connection enum
	ResourcePool_ResourceType_Ethernet_Connection ResourcePool_ResourceType = 33
	// DMTF_reserved enum
	ResourcePool_ResourceType_DMTF_reserved ResourcePool_ResourceType = 34
	// Vendor_Reserved enum
	ResourcePool_ResourceType_Vendor_Reserved ResourcePool_ResourceType = 35
)

// GetResourceAllocationSettingData returns the ResourceAllocationSettingData for the given role and range
func (rp *ResourcePool) GetResourceAllocationSettingData(crole SettingsDefineCapabilities_ValueRole, crange SettingsDefineCapabilities_ValueRange) (setting *ResourceAllocationSettingData, err error) {
	var (
		allocationCapabilities        []*wmiext.Instance
		settingsDefineCapabilities    []*wmiext.Instance
		settingsDefineCapability      SettingsDefineCapabilities
		resourceAllocationSettingData ResourceAllocationSettingData
	)

	if allocationCapabilities, err = rp.GetAllRelated("Msvm_AllocationCapabilities"); err != nil {
		return
	}

	for _, allocationCapability := range allocationCapabilities {
		if settingsDefineCapabilities, err = allocationCapability.GetReferences(Msvm_SettingsDefineCapabilities); err != nil {
			return
		}

		for _, settingsDefineCapabilityInst := range settingsDefineCapabilities {
			settingsDefineCapability = SettingsDefineCapabilities{}

			if err = settingsDefineCapabilityInst.GetAll(&settingsDefineCapability); err != nil {
				return
			}

			valueRange := int(settingsDefineCapability.ValueRange)
			valueRole := int(settingsDefineCapability.ValueRole)
			if valueRange != int(crange) || valueRole != int(crole) {
				continue
			}
			// Found the match
			partComponent := settingsDefineCapability.PartComponent

			resourceAllocationSettingData = ResourceAllocationSettingData{}

			if err = rp.GetService().GetObjectAsObject(partComponent, &resourceAllocationSettingData); err != nil {
				return
			}

			return &resourceAllocationSettingData, nil
		}
	}
	return nil, errors.Wrapf(wmiext.NotFound, "GetResourceAllocationSettingData [%d] [%d]", crole, crange)
}

// GetDefaultResourceAllocationSettingData returns the default ResourceAllocationSettingData
func (rp *ResourcePool) GetDefaultResourceAllocationSettingData() (*ResourceAllocationSettingData, error) {
	return rp.GetResourceAllocationSettingData(SettingsDefineCapabilities_ValueRole_Default, SettingsDefineCapabilities_ValueRange_Point)
}

func GetPrimordialResourcePool(session *wmiext.Service, rtype ResourcePool_ResourceType) (*ResourcePool, error) {
	var (
		err error
	)
	rptype := GetResourceTypeValue(rtype)
	col, err := GetResourcePools[any](session, true, rptype)
	if err != nil {
		return nil, err
	}
	if len(col) == 0 {
		return nil, errors.Wrapf(wmiext.NotFound, "ResourcePool [%s]", rptype.String())
	}
	inst, err := col[0].CloneInstance()
	if err != nil {
		return nil, err
	}
	rp := ResourcePool{}
	return &rp, inst.GetAll(&rp)
}

func GetResourcePools[T any](session *wmiext.Service, isPrimordial bool, resType *ResourceTypeValue) (col []*wmiext.Instance, err error) {
	wquery := "SELECT * FROM Msvm_ResourcePool"
	if isPrimordial {
		wquery += " WHERE Primordial=True"
	}
	if resType != nil {
		if resType.ResourceType > 0 {
			if isPrimordial {
				wquery += fmt.Sprintf(" AND ResourceType=%d", resType.ResourceType)
			} else {
				wquery += fmt.Sprintf(" WHERE ResourceType=%d", resType.ResourceType)
			}
		}
		if len(resType.ResourceSubType) > 0 {
			wquery += fmt.Sprintf(" AND ResourceSubType='%s'", resType.ResourceSubType)
		}
		if len(resType.OtherResourceType) > 0 {
			wquery += fmt.Sprintf(" AND OtherResourceType='%s'", resType.OtherResourceType)
		}
	}
	if col, err = session.FindInstances(wquery); err != nil {
		return nil, err
	}
	if len(col) == 0 {
		err = errors.Wrapf(wmiext.NotFound, "Cim_ResourcePool Primordial[%s] Type[%s]", isPrimordial, resType.String())
	}
	return
}
