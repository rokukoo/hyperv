package disk

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/allocation"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualHardDisk struct {
	*allocation.StorageAllocationSettingData
}

func NewVirtualHardDisk(instance *wmiext.Instance) (*VirtualHardDisk, error) {
	storageAllocationSettingData := &allocation.StorageAllocationSettingData{}
	if err := instance.GetAll(storageAllocationSettingData); err != nil {
		return nil, err
	}
	return &VirtualHardDisk{storageAllocationSettingData}, nil
}

func (vhd *VirtualHardDisk) GetPath() string {
	return vhd.HostResource[0]
}

func (vhd *VirtualHardDisk) GetDrive() (resourceAllocationSettingData *resource.ResourceAllocationSettingData, err error) {
	resourceAllocationSettingData = &resource.ResourceAllocationSettingData{}
	parent := vhd.Parent
	if err = vhd.GetService().GetObjectAsObject(parent, resourceAllocationSettingData); err != nil {
		return nil, err
	}
	return resourceAllocationSettingData, nil
}

func GetVirtualHardDiskByPath(session *wmiext.Service, path string) (virtualHardDisk *VirtualHardDisk, err error) {
	wquery := fmt.Sprintf("SELECT * FROM Msvm_StorageAllocationSettingData")
	instances, err := session.FindInstances(wquery)
	if err != nil {
		return nil, err
	}
	var hostResourceProp any
	for _, instance := range instances {
		hostResourceProp, _, _, err = instance.GetAsAny("HostResource")
		if err != nil {
			return nil, err
		}
		if hostResourceProp == nil {
			continue
		}
		hostResource := hostResourceProp.([]interface{})
		if len(hostResource) == 0 {
			continue
		}
		if hostResource[0] != path {
			continue
		}
		if virtualHardDisk, err = NewVirtualHardDisk(instance); err != nil {
			return nil, err
		}
		return virtualHardDisk, nil
	}
	return nil, wmiext.NotFound
}
