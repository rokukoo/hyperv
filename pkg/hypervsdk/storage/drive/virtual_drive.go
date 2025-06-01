package drive

import (
	"github.com/rokukoo/hyperv/pkg/hypervsdk/resource"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

type VirtualDrive struct {
	*resource.ResourceAllocationSettingData
}

// NewVirtualDrive creates a new VirtualDrive instance
func NewVirtualDrive(instance *wmiext.Instance) (*VirtualDrive, error) {
	resourceAllocationSettingData := resource.ResourceAllocationSettingData{}
	if err := instance.GetAll(&resourceAllocationSettingData); err != nil {
		return nil, err
	}
	return &VirtualDrive{&resourceAllocationSettingData}, nil
}

func (vdrive *VirtualDrive) GetController() (*resource.ResourceAllocationSettingData, error) {
	parent := vdrive.Parent

	resourceAllocationSettingData := resource.ResourceAllocationSettingData{}
	if err := vdrive.GetService().GetObjectAsObject(parent, &resourceAllocationSettingData); err != nil {
		return nil, err
	}
	return &resourceAllocationSettingData, nil
}

func (vdrive *VirtualDrive) GetControllerLocation() string {
	return vdrive.AddressOnParent
}

func (vdrive *VirtualDrive) GetControllerNumber() (string, error) {
	controller, err := vdrive.GetController()
	if err != nil {
		return "0", err
	}
	val := controller.Address
	if len(val) == 0 {
		return "0", nil
	}
	return val, err
}
