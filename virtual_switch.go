package hypervctl

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking/networking_service"
	hypervsdk "github.com/rokukoo/hypervctl/pkg/hypervsdk/utils"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualSwitchType = int

const (
	VirtualSwitchTypePrivate VirtualSwitchType = iota
	VirtualSwitchTypeInternal
	VirtualSwitchTypeExternalBridge
	VirtualSwitchTypeExternalDirect
)

// VirtualSwitch represents a Hyper-V virtual switch
type VirtualSwitch struct {
	// Name of the virtual switch
	Name string `json:"name"`
	// Description of the virtual switch
	Description string `json:"description"`
	// Type of the virtual switch, can be private, internal, external bridge or external direct
	Type                  VirtualSwitchType `json:"type"`
	PhysicalAdapter       *string           `json:"physical_adapter"`
	virtualEthernetSwitch *networking.VirtualEthernetSwitch
}

func NewVirtualSwitch(virtualEthernetSwitch *networking.VirtualEthernetSwitch) (*VirtualSwitch, error) {
	vsw := &VirtualSwitch{}
	vsw.virtualEthernetSwitch = virtualEthernetSwitch
	return vsw, vsw.update(virtualEthernetSwitch)
}

func (vsw *VirtualSwitch) update(virtualEthernetSwitch *networking.VirtualEthernetSwitch) (err error) {
	vsw.virtualEthernetSwitch = virtualEthernetSwitch
	vsw.Name = virtualEthernetSwitch.ElementName
	virtualEthernetSwitchSettingData, err := virtualEthernetSwitch.ActiveVirtualEthernetSwitchSettingData()
	if err != nil {
		return
	}
	vsw.Description = virtualEthernetSwitchSettingData.Notes[0]
	vsw.Type, err = vsw.GetType()
	if err != nil {
		return errors.Wrap(err, "failed to get virtual switch type")
	}
	externalPortAllocSettings, err := virtualEthernetSwitch.GetExternalPortAllocSettings()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return errors.Wrap(err, "failed to get external port allocation setting data")
	}
	if externalPortAllocSettings != nil {
		externalPort := &networking.ExternalEthernetPort{}
		if err = virtualEthernetSwitch.GetService().GetObjectAsObject(externalPortAllocSettings.HostResource[0], externalPort); err != nil {
			return errors.Wrap(err, "failed to find external ethernet port")
		}
		vsw.PhysicalAdapter = &externalPort.Name
	}
	return nil
}

func (vsw *VirtualSwitch) GetType() (VirtualSwitchType, error) {
	var (
		err error
		internalPortAllocSettings,
		externalPortAllocSettings *networking.EthernetPortAllocationSettingData
	)
	if internalPortAllocSettings, err = vsw.virtualEthernetSwitch.GetInternalPortAllocSettings(); err != nil && !errors.Is(err, wmiext.NotFound) {
		return 0, errors.Wrap(err, "failed to get internal port allocation setting data")
	}
	if externalPortAllocSettings, err = vsw.virtualEthernetSwitch.GetExternalPortAllocSettings(); err != nil && !errors.Is(err, wmiext.NotFound) {
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

func (vsw *VirtualSwitch) ChangeType(switchType VirtualSwitchType, adapter *string) (err error) {
	virtualSwitch := vsw.virtualEthernetSwitch
	var vsms *networking_service.VirtualEthernetSwitchManagementService

	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return errors.Wrap(err, "failed to get virtual switch management service")
	}

	switch switchType {
	case VirtualSwitchTypePrivate:
		// Change virtual switch to private
		if err = vsms.ClearInternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear internal port allocation setting data")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		return nil
	case VirtualSwitchTypeInternal:
		// Change virtual switch to internal
		var resourceSettings []*networking.EthernetPortAllocationSettingData
		_, err = virtualSwitch.GetInternalPortAllocSettings()
		if err != nil {
			if errors.Is(err, wmiext.NotFound) {
				epads, err := vsms.DefaultInternalPortAllocationSettingData(virtualSwitch.ElementName)
				if err != nil {
					return errors.Wrap(err, "failed to get internal port allocation settings")
				}
				resourceSettings = append(resourceSettings, epads)
			} else {
				return errors.Wrap(err, "failed to get internal port allocation settings")
			}
		}
		settingData, err := virtualSwitch.ActiveVirtualEthernetSwitchSettingData()
		if err != nil {
			return errors.Wrap(err, "failed to get active virtual ethernet switch setting data")
		}
		if _, err = vsms.AddResourceSettings(settingData, resourceSettings); err != nil {
			return errors.Wrap(err, "failed to add allocation settings")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		return nil
	case VirtualSwitchTypeExternalBridge:
		if adapter == nil {
			return errors.New("adapter is required for external bridge")
		}
		var resourceSettings []*networking.EthernetPortAllocationSettingData
		_, err = virtualSwitch.GetInternalPortAllocSettings()
		if err != nil {
			if errors.Is(err, wmiext.NotFound) {
				epads, err := vsms.DefaultInternalPortAllocationSettingData(virtualSwitch.ElementName)
				if err != nil {
					return errors.Wrap(err, "failed to get internal port allocation settings")
				}
				resourceSettings = append(resourceSettings, epads)
			} else {
				return errors.Wrap(err, "failed to get internal port allocation settings")
			}
		}
		if err = vsms.ClearExternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		epads, err := vsms.DefaultExternalPortAllocationSettingData(virtualSwitch.ElementName, []string{*adapter})
		if err != nil {
			return errors.Wrap(err, "failed to get external port allocation settings")
		}
		resourceSettings = append(resourceSettings, epads)
		settingData, err := virtualSwitch.ActiveVirtualEthernetSwitchSettingData()
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
		if err = vsms.ClearInternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear internal port allocation setting data")
		}
		if err = vsms.ClearExternalPortAllocationSettingData(virtualSwitch); err != nil {
			return errors.Wrap(err, "failed to clear external port allocation setting data")
		}
		epads, err := vsms.DefaultExternalPortAllocationSettingData(virtualSwitch.ElementName, []string{*adapter})
		if err != nil {
			return errors.Wrap(err, "failed to get external port allocation settings")
		}
		resourceSettings = append(resourceSettings, epads)
		settingData, err := virtualSwitch.ActiveVirtualEthernetSwitchSettingData()
		if err != nil {
			return errors.Wrap(err, "failed to get active virtual ethernet switch setting data")
		}
		if _, err = vsms.AddResourceSettings(settingData, resourceSettings); err != nil {
			return errors.Wrap(err, "failed to add allocation settings")
		}
		return nil
	default:
		return errors.New("invalid virtual switch type")
	}
}

func CreateVirtualSwitch(name string, switchType VirtualSwitchType, adapter *string) (*VirtualSwitch, error) {
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
		return NewVirtualSwitch(vsw)
	case VirtualSwitchTypeInternal:
		// Build internal virtual switch
		if vsw, err = CreateInternalVirtualSwitch(name); err != nil {
			return nil, errors.Wrap(err, "failed to create internal virtual")
		}
		return NewVirtualSwitch(vsw)
	case VirtualSwitchTypeExternalBridge:
		// Build external virtual switch
		if vsw, err = CreateExternalVirtualSwitch(name, *adapter, true); err != nil {
			return nil, errors.Wrap(err, "failed to create external virtual switch")
		}
		return NewVirtualSwitch(vsw)
	case VirtualSwitchTypeExternalDirect:
		// Build external virtual switch directly
		if vsw, err = CreateExternalVirtualSwitch(name, *adapter, false); err != nil {
			return nil, errors.Wrap(err, "failed to create external virtual switch")
		}
		return NewVirtualSwitch(vsw)
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
	vSwitch, err := vsms.CreateExternalVirtualSwitch(networkInterfaceDescription, portName, portName, switchSettingData, internalport)
	if err != nil {
		return nil, err
	}
	return vSwitch, nil
}

// FirstVirtualSwitchByName returns the first virtual switch by name
func FirstVirtualSwitchByName(name string) (*VirtualSwitch, error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		err  error
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return nil, err
	}
	virtualSwitch, err := vsms.FirstVirtualSwitchByName(name)
	if err != nil {
		return nil, err
	}
	return NewVirtualSwitch(virtualSwitch)
}

func MustFirstVirtualSwitchByName(name string) *VirtualSwitch {
	vsw, err := FirstVirtualSwitchByName(name)
	if err != nil {
		panic(err)
	}
	return vsw
}

func GetVirtualSwitchTypeByName(name string) (VirtualSwitchType, error) {
	vsw, err := FirstVirtualSwitchByName(name)
	if err != nil {
		return 0, err
	}
	return vsw.Type, nil
}

func ChangeVirtualSwitchTypeByName(name string, switchType VirtualSwitchType, adapter *string) error {
	vsw, err := FirstVirtualSwitchByName(name)
	if err != nil {
		return err
	}
	return vsw.ChangeType(switchType, adapter)
}

// DeleteVirtualSwitchByName removes a virtual switch by name
func DeleteVirtualSwitchByName(name string) (err error) {
	var (
		vsms *networking_service.VirtualEthernetSwitchManagementService
		vsw  *networking.VirtualEthernetSwitch
	)
	if vsms, err = networking_service.LocalVirtualEthernetSwitchManagementService(); err != nil {
		return
	}
	if vsw, err = vsms.FirstVirtualSwitchByName(name); err != nil {
		return
	}
	if err = vsms.DestroySystem(vsw); err != nil {
		return
	}
	return
}

func ListAvailablePhysicalNetworkAdapters() ([]string, error) {
	var nics []string
	service, err := hypervsdk.NewLocalHyperVService()
	if err != nil {
		return nil, err
	}
	externalEthernetPorts, err := networking.ListEnabledExternalEthernetPort(service)
	if err != nil {
		return nil, err
	}
	for _, port := range externalEthernetPorts {
		nics = append(nics, port.Name)
	}

	wiFiPorts, err := networking.ListEnabledWiFiPort(service)
	if err != nil {
		return nil, err
	}
	for _, port := range wiFiPorts {
		nics = append(nics, port.Name)
	}

	return nics, nil
}
