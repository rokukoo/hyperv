package hypervctl

import (
	"github.com/google/uuid"
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/network/virtualswitch"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/wmictl"
)

type VirtualSwitchType = int

const (
	VirtualSwitchTypePrivate VirtualSwitchType = iota
	VirtualSwitchTypeInternal
	VirtualSwitchTypeExternal
	VirtualSwitchTypeExternalDirect
)

type VirtualSwitch struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Status      int               `json:"status"`
	Type        VirtualSwitchType `json:"type"`
	*virtualswitch.VirtualSwitch
}

func virtualSwitch(vsw *virtualswitch.VirtualSwitch) (*VirtualSwitch, error) {
	id := vsw.ID()
	name, err := vsw.GetPropertyElementName()
	if err != nil {
		return nil, err
	}
	description, err := vsw.GetPropertyDescription()
	if err != nil {
		return nil, err
	}
	enabledState, err := vsw.GetPropertyEnabledState()
	if err != nil {
		return nil, err
	}
	return &VirtualSwitch{
		Id:            id,
		Name:          name,
		Description:   description,
		Status:        int(enabledState),
		VirtualSwitch: vsw,
	}, nil
}

func FindVirtualSwitchByName(name string) (*VirtualSwitch, error) {
	whost := host.NewWmiLocalHost()
	vswitch, err := virtualswitch.GetVirtualSwitch(whost, name)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get virtual switch")
	}
	vSwitch, err := virtualSwitch(vswitch)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create virtual switch")
	}
	return vSwitch, nil
}

// CreatePrivateVirtualSwitch creates a private virtual switch
func CreatePrivateVirtualSwitch(name string) (*VirtualSwitch, error) {
	var (
		vsms *wmictl.VirtualEthernetSwitchManagementService
		err  error
	)
	if vsms, err = wmictl.NewLocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	setting, err := virtualswitch.GetVirtualEthernetSwitchSettingData(vsms.GetWmiHost(), "private")
	if err != nil {
		return nil, err
	}
	if err := setting.SetPropertyElementName(name); err != nil {
		return nil, err
	}
	vswitch, err := vsms.CreatePrivateVirtualSwitch(setting)
	if err != nil {
		return nil, err
	}
	vsw, err := virtualSwitch(vswitch)
	if err != nil {
		return nil, err
	}
	return vsw, nil
}

func CreateInternalVirtualSwitch(name string) (*VirtualSwitch, error) {
	var (
		vsms *wmictl.VirtualEthernetSwitchManagementService
		err  error
	)
	if vsms, err = wmictl.NewLocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	setting, err := virtualswitch.GetVirtualEthernetSwitchSettingData(vsms.GetWmiHost(), "internal")
	if err != nil {
		return nil, err
	}
	if err = setting.SetPropertyElementName(name); err != nil {
		return nil, err
	}
	vswitch, err := vsms.CreateInternalVirtualSwitch(name, setting)
	if err != nil {
		return nil, err
	}
	vsw, err := virtualSwitch(vswitch)
	if err != nil {
		return nil, err
	}
	return vsw, nil
}

func CreateExternalVirtualSwitch(name, networkInterfaceName string, internalport bool) (*VirtualSwitch, error) {
	var (
		vsms *wmictl.VirtualEthernetSwitchManagementService
		err  error
	)

	if vsms, err = wmictl.NewLocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	setting, err := virtualswitch.GetVirtualEthernetSwitchSettingData(vsms.GetWmiHost(), "external")
	if err != nil {
		return nil, err
	}
	if err = setting.SetPropertyElementName(name); err != nil {
		return nil, err
	}
	portName := uuid.NewString()
	netAdapter, err := wmictl.FindNetAdapterByInterfaceDescription(networkInterfaceName)
	if err != nil {
		return nil, err
	}
	physicalNicName := netAdapter.Name

	vswitch, err := vsms.CreateExternalVirtualSwitch(physicalNicName, portName, portName, setting, internalport)
	if err != nil {
		return nil, err
	}
	vsw, err := virtualSwitch(vswitch)
	if err != nil {
		return nil, err
	}
	return vsw, nil
}

func (vsw *VirtualSwitch) Create() error {
	vswType := vsw.Type
	switch vswType {
	case 0:
		// Create private virtual switch
		_, err := CreatePrivateVirtualSwitch(vsw.Name)
		if err != nil {
			return errors.Wrap(err, "failed to create private virtual switch")
		}
		break
	case 1:
		// Create internal virtual switch
		_, err := CreateInternalVirtualSwitch(vsw.Name)
		if err != nil {
			return errors.Wrap(err, "failed to create internal virtual")
		}
		break
	case 2:
		// Create external virtual switch
		break
	case 3:
		// Create external virtual switch directly
		break
	}
	return nil
}

// DeleteVirtualSwitchByName removes a virtual switch by name
func DeleteVirtualSwitchByName(name string) (bool, error) {
	var (
		vsms *wmictl.VirtualEthernetSwitchManagementService
		err  error
	)
	if vsms, err = wmictl.NewLocalVirtualEthernetSwitchManagementService(); err != nil {
		return false, err
	}
	vsw, err := vsms.FindVirtualSwitchByName(name)
	if err != nil {
		return false, err
	}
	if err = vsms.DeleteVirtualSwitch(vsw); err != nil {
		return false, err
	}
	return true, nil
}
