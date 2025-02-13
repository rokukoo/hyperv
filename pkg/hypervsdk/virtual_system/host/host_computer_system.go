package host

import (
	"fmt"
	"github.com/microsoft/wmi/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking/switch_extension"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type HostComputerSystem struct {
	*virtual_system.ComputerSystem
}

// GetHostComputerSystem gets an existing virtual machine
func GetHostComputerSystem() (*HostComputerSystem, error) {
	var (
		vm  = &virtual_system.ComputerSystem{}
		err error
	)
	// TODO: Add a filter for the host computer system
	wquery := fmt.Sprintf("SELECT * FROM %s WHERE Description != 'Microsoft Virtual Machine' AND Description != 'Microsoft 虚拟机'", virtual_system.Msvm_ComputerSystem)
	s, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	if err = s.Session.FindFirstObject(wquery, vm); err != nil {
		return nil, err
	}
	return &HostComputerSystem{vm}, nil
}

func (hc *HostComputerSystem) GetInstalledEthernetSwitchExtensions() (col []*switch_extension.InstalledEthernetSwitchExtension, err error) {
	var installedEthernetSwitchExtension *switch_extension.InstalledEthernetSwitchExtension
	installedEthernetSwitchExtensions, err := hc.GetAllRelated("Msvm_InstalledEthernetSwitchExtension")
	if err != nil {
		return
	}

	for _, inst := range installedEthernetSwitchExtensions {
		if installedEthernetSwitchExtension, err = switch_extension.NewInstalledEthernetSwitchExtension(inst); err != nil {
			return
		}
		col = append(col, installedEthernetSwitchExtension)
	}

	return
}

func (hc *HostComputerSystem) GetFeatureCapability(featureName string) (*switch_extension.EthernetSwitchFeatureCapabilities, error) {
	installedEthernetSwitchExtensions, err := hc.GetInstalledEthernetSwitchExtensions()
	if err != nil {
		return nil, err
	}

	for _, ext := range installedEthernetSwitchExtensions {
		fc, err := ext.GetFeatureCapabilityByName(featureName)
		if err != nil {
			if errors.IsNotFound(err) {
				continue
			}
			return nil, err
		}
		// If found, return
		return fc, nil
	}

	return nil, errors.Wrapf(errors.NotFound, "FeatureCapability [%s]", featureName)
}

func (hc *HostComputerSystem) GetDefaultPortSettingData(featureName, className string) (*wmiext.Instance, error) {
	capability, err := hc.GetFeatureCapability(featureName)
	if err != nil {
		return nil, err
	}
	defer capability.Close()

	return capability.GetRelated(className)
}

// DefaultEthernetSwitchPortBandwidthSettingData returns the default EthernetSwitchPortBandwidthSettingData
func DefaultEthernetSwitchPortBandwidthSettingData() (*switch_extension.EthernetSwitchPortBandwidthSettingData, error) {
	hc, err := GetHostComputerSystem()
	if err != nil {
		return nil, err
	}
	inst, err := hc.GetDefaultPortSettingData("Ethernet Switch Port Bandwidth Settings", "Msvm_EthernetSwitchPortBandwidthSettingData")
	if err != nil {
		return nil, err
	}
	spbs, err := switch_extension.NewEthernetSwitchPortBandwidthSettingData(inst)
	if err != nil {
		return nil, err
	}
	return spbs, nil
}
