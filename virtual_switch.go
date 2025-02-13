package hypervctl

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking/networking_service"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualSwitchType = int

const (
	VirtualSwitchTypePrivate VirtualSwitchType = iota
	VirtualSwitchTypeInternal
	VirtualSwitchTypeExternalBridge
	VirtualSwitchTypeExternalDirect
)

type VirtualSwitch struct {
	*networking.VirtualEthernetSwitch
}

func NewVirtualSwitch(virtualEthernetSwitch *networking.VirtualEthernetSwitch) (*VirtualSwitch, error) {
	vsw := &VirtualSwitch{}
	vsw.VirtualEthernetSwitch = virtualEthernetSwitch
	return vsw, vsw.update(virtualEthernetSwitch)
}

func (vsw *VirtualSwitch) update(virtualEthernetSwitch *networking.VirtualEthernetSwitch) error {
	vsw.VirtualEthernetSwitch = virtualEthernetSwitch
	return nil
}

func GetVirtualSwitchTypeByName(name string) (VirtualSwitchType, error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		vsw  *networking.VirtualEthernetSwitch
		err  error
		internalPortAllocSettings,
		externalPortAllocSettings *networking.EthernetPortAllocationSettingData
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return 0, errors.Wrap(err, "failed to get virtual switch management service")
	}
	if vsw, err = vsms.FirstVirtualSwitchByName(name); err != nil {
		return 0, errors.Wrap(err, "failed to find virtual switch")
	}
	if internalPortAllocSettings, err = vsw.GetInternalPortAllocSettings(); err != nil && !errors.Is(err, wmiext.NotFound) {
		return 0, errors.Wrap(err, "failed to get internal port allocation setting data")
	}
	if externalPortAllocSettings, err = vsw.GetExternalPortAllocSettings(); err != nil && !errors.Is(err, wmiext.NotFound) {
		return 0, errors.Wrap(err, "failed to get internal port allocation setting data")
	}
	if internalPortAllocSettings != nil {
		if externalPortAllocSettings != nil {
			return VirtualSwitchTypeExternalBridge, nil
		} else {
			return VirtualSwitchTypeInternal, nil
		}
	} else {
		if externalPortAllocSettings != nil {
			return VirtualSwitchTypeExternalDirect, nil
		} else {
			return VirtualSwitchTypePrivate, nil
		}
	}
}

func ChangeVirtualSwitchTypeByName(name string, switchType VirtualSwitchType, adapter *string) error {
	var vsms *networking_service.VirtualEthernetSwitchManagementService
	var err error
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return errors.Wrap(err, "failed to get virtual switch management service")
	}

	var vsw *networking.VirtualEthernetSwitch
	if vsw, err = vsms.FirstVirtualSwitchByName(name); err != nil {
		return errors.Wrap(err, "failed to find virtual switch")
	}
	switch switchType {
	case VirtualSwitchTypePrivate:
		// Change virtual switch to private
		if err = vsms.ClearInternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear internal port allocation setting data")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		return nil
	case VirtualSwitchTypeInternal:
		// Change virtual switch to internal
		var resourceSettings []*networking.EthernetPortAllocationSettingData
		_, err = vsw.GetInternalPortAllocSettings()
		if err != nil {
			if errors.Is(err, wmiext.NotFound) {
				epads, err := vsms.DefaultInternalPortAllocationSettingData(vsw.ElementName)
				if err != nil {
					return errors.Wrap(err, "failed to get internal port allocation settings")
				}
				resourceSettings = append(resourceSettings, epads)
			} else {
				return errors.Wrap(err, "failed to get internal port allocation settings")
			}
		}
		settingData, err := vsw.ActiveVirtualEthernetSwitchSettingData()
		if err != nil {
			return errors.Wrap(err, "failed to get active virtual ethernet switch setting data")
		}
		if _, err = vsms.AddResourceSettings(settingData, resourceSettings); err != nil {
			return errors.Wrap(err, "failed to add allocation settings")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		return nil
	case VirtualSwitchTypeExternalBridge:
		if adapter == nil {
			return errors.New("adapter is required for external bridge")
		}
		var resourceSettings []*networking.EthernetPortAllocationSettingData
		_, err = vsw.GetInternalPortAllocSettings()
		if err != nil {
			if errors.Is(err, wmiext.NotFound) {
				epads, err := vsms.DefaultInternalPortAllocationSettingData(vsw.ElementName)
				if err != nil {
					return errors.Wrap(err, "failed to get internal port allocation settings")
				}
				resourceSettings = append(resourceSettings, epads)
			} else {
				return errors.Wrap(err, "failed to get internal port allocation settings")
			}
		}
		if err = vsms.ClearExternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		epads, err := vsms.DefaultExternalPortAllocationSettingData(vsw.ElementName, []string{*adapter})
		if err != nil {
			return errors.Wrap(err, "failed to get external port allocation settings")
		}
		resourceSettings = append(resourceSettings, epads)
		settingData, err := vsw.ActiveVirtualEthernetSwitchSettingData()
		if err != nil {
			return errors.Wrap(err, "failed to get active virtual ethernet switch setting data")
		}
		if _, err = vsms.AddResourceSettings(settingData, resourceSettings); err != nil {
			return errors.Wrap(err, "failed to add allocation settings")
		}
		return nil
	// Change virtual switch to external bridge
	case VirtualSwitchTypeExternalDirect:
		if adapter == nil {
			return errors.New("adapter is required for external direct")
		}
		var resourceSettings []*networking.EthernetPortAllocationSettingData
		// Change virtual switch to private
		if err = vsms.ClearInternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear internal port allocation setting data")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(vsw); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		epads, err := vsms.DefaultExternalPortAllocationSettingData(vsw.ElementName, []string{*adapter})
		if err != nil {
			return errors.Wrap(err, "failed to get external port allocation settings")
		}
		resourceSettings = append(resourceSettings, epads)
		settingData, err := vsw.ActiveVirtualEthernetSwitchSettingData()
		if err != nil {
			return errors.Wrap(err, "failed to get active virtual ethernet switch setting data")
		}
		if _, err = vsms.AddResourceSettings(settingData, resourceSettings); err != nil {
			return errors.Wrap(err, "failed to add allocation settings")
		}
		return nil
	}
	return nil
}

func CreateVirtualSwitch(name string, description string, switchType VirtualSwitchType, adapter *string) (*networking.VirtualEthernetSwitch, error) {
	var (
		vsw = &networking.VirtualEthernetSwitch{}
		err error
	)
	switch switchType {
	case VirtualSwitchTypePrivate:
		// Build private virtual switch
		if vsw, err = CreatePrivateVirtualSwitch(name); err != nil {
			return nil, errors.Wrap(err, "failed to create private virtual switch")
		}
		return vsw, nil
	case VirtualSwitchTypeInternal:
		// Build internal virtual switch
		if vsw, err = CreateInternalVirtualSwitch(name); err != nil {
			return nil, errors.Wrap(err, "failed to create internal virtual")
		}
		return vsw, nil
	case VirtualSwitchTypeExternalBridge:
		// Build external virtual switch
		if vsw, err = CreateExternalVirtualSwitch(name, *adapter, true); err != nil {
			return nil, errors.Wrap(err, "failed to create external virtual switch")
		}
		return vsw, nil
	case VirtualSwitchTypeExternalDirect:
		// Build external virtual switch directly
		if vsw, err = CreateExternalVirtualSwitch(name, *adapter, false); err != nil {
			return nil, errors.Wrap(err, "failed to create external virtual switch")
		}
		return vsw, nil
	default:
		return nil, errors.New("invalid virtual switch type")
	}
}

// CreatePrivateVirtualSwitch creates a private virtual switch
func CreatePrivateVirtualSwitch(name string) (*networking.VirtualEthernetSwitch, error) {
	var (
		vsms    *networking_service.VirtualEthernetSwitchManagementService
		setting *networking.VirtualEthernetSwitchSettingData
		err     error
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	if setting, err = vsms.GetVirtualEthernetSwitchSettingData(name); err != nil {
		return nil, err
	}
	return vsms.CreatePrivateVirtualSwitch(setting)
}

func CreateInternalVirtualSwitch(name string) (*networking.VirtualEthernetSwitch, error) {
	var (
		vsms    *networking_service.VirtualEthernetSwitchManagementService
		setting *networking.VirtualEthernetSwitchSettingData
		err     error
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	if setting, err = vsms.GetVirtualEthernetSwitchSettingData(name); err != nil {
		return nil, err
	}
	return vsms.CreateInternalVirtualSwitch(name, setting)
}

func CreateExternalVirtualSwitch(name, networkInterfaceDescription string, internalport bool) (*networking.VirtualEthernetSwitch, error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		err  error
	)

	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	switchSettingData, err := vsms.GetVirtualEthernetSwitchSettingData(name)
	if err != nil {
		return nil, err
	}
	portName := uuid.NewString()
	vswitch, err := vsms.CreateExternalVirtualSwitch(networkInterfaceDescription, portName, portName, switchSettingData, internalport)
	if err != nil {
		return nil, err
	}
	return vswitch, nil
}

// DeleteVirtualSwitchByName removes a virtual switch by name
func DeleteVirtualSwitchByName(name string) (bool, error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		vsw  *networking.VirtualEthernetSwitch
		err  error
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return false, err
	}
	if vsw, err = vsms.FirstVirtualSwitchByName(name); err != nil {
		return false, err
	}
	if err = vsms.DestroySystem(vsw); err != nil {
		return false, err
	}
	return true, nil
}

func FindVirtualSwitchByName(name string) (*networking.VirtualEthernetSwitch, error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		err  error
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	return vsms.FirstVirtualSwitchByName(name)
}
