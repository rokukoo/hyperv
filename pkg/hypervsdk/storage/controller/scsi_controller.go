package controller

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/drive"
)

type SCSIControllerSettings struct {
	*resource.ResourceAllocationSettingData
}

// NewSCSIControllerSettings creates a new SCSIControllerSettings instance
func NewSCSIControllerSettings(resourceAllocationSettingData *resource.ResourceAllocationSettingData) *SCSIControllerSettings {
	return &SCSIControllerSettings{resourceAllocationSettingData}
}

func checkIfLocationExists(resourceAllocationSettingData []*resource.ResourceAllocationSettingData, locationNumber int) (bool, error) {
	for _, inst := range resourceAllocationSettingData {
		syntheticDiskDrive := &drive.SyntheticDiskDrive{VirtualDrive: &drive.VirtualDrive{ResourceAllocationSettingData: inst}}
		loc := syntheticDiskDrive.GetControllerLocation()
		if loc == fmt.Sprintf("%d", locationNumber) {
			return true, nil
		}

	}
	return false, nil
}

func (settings *SCSIControllerSettings) GetFreeLocation() (int32, error) {
	// Get all drives - Is this
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

func (settings *SCSIControllerSettings) getResourceAllocationSettingData(rtype resource.ResourceAllocationSettingData_ResourceType) (col []*resource.ResourceAllocationSettingData, err error) {
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
