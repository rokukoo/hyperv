package controller

import (
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
)

type IDEControllerSettings struct {
	*resource.ResourceAllocationSettingData
}

// NewIDEControllerSettings creates a new IDEControllerSettings instance
func NewIDEControllerSettings(resourceAllocationSettingData *resource.ResourceAllocationSettingData) *IDEControllerSettings {
	return &IDEControllerSettings{resourceAllocationSettingData}
}

// GetFreeLocation returns the first free location for a new drive
func (settings *IDEControllerSettings) GetFreeLocation() (int32, error) {
	// Get all drives
	resourceAllocationSettingDatas, err := settings.getResourceAllocationSettingData(resource.ResourceAllocationSettingData_ResourceType_Disk_Drive)
	if err != nil {
		return -1, err
	}
	dvdDriveResourceAllocationSettingDatas, err := settings.getResourceAllocationSettingData(resource.ResourceAllocationSettingData_ResourceType_DVD_drive)
	if err != nil {
		return -1, err
	}
	resourceAllocationSettingDatas = append(resourceAllocationSettingDatas, dvdDriveResourceAllocationSettingDatas...)

	freeLocation := 0
	var exists bool
	for range resourceAllocationSettingDatas {
		if exists, err = checkIfLocationExists(resourceAllocationSettingDatas, freeLocation); err != nil {
			return -1, err
		}
		if exists {
			freeLocation = freeLocation + 1
			continue
		}
		break
	}
	return int32(freeLocation), nil
}

func (settings *IDEControllerSettings) getResourceAllocationSettingData(rtype resource.ResourceAllocationSettingData_ResourceType) (col []*resource.ResourceAllocationSettingData, err error) {
	var (
		resourceAllocationSettingData *resource.ResourceAllocationSettingData
	)
	resourceType := uint16(rtype)
	resourceAllocationSettingDatas, err := settings.GetAllRelated("Msvm_ResourceAllocationSettingData")
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
