package network_adapter

import (
	"fmt"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/networking"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

type VirtualNetworkAdapter struct {
	*networking.EthernetPortAllocationSettingData
}

// NewVirtualNetworkAdapterFromInstance creates a new VirtualNetworkAdapter instance from a WMI instance
func NewVirtualNetworkAdapterFromInstance(instance *wmiext.Instance) (*VirtualNetworkAdapter, error) {
	epasd := &networking.EthernetPortAllocationSettingData{}
	if err := instance.GetAll(epasd); err != nil {
		return nil, err
	}
	return &VirtualNetworkAdapter{epasd}, nil
}

func (vna *VirtualNetworkAdapter) Clone() (*VirtualNetworkAdapter, error) {
	instance, err := vna.CloneInstance()
	if err != nil {
		return nil, err
	}
	return NewVirtualNetworkAdapterFromInstance(instance)
}

func (vna *VirtualNetworkAdapter) GetEthernetPortAllocationSettingData() (epas *networking.EthernetPortAllocationSettingData, err error) {
	instance, err := vna.GetRelated("Msvm_EthernetPortAllocationSettingData")
	if err != nil {
		return
	}
	return networking.NewEthernetPortAllocationSettingDataFromInstance(instance)
}

func (vna *VirtualNetworkAdapter) SetElementName(elementName string) error {
	vna.ElementName = elementName
	return vna.Put("ElementName", elementName)
}

func FirstVirtualNetworkAdapterByName(session *wmiext.Service, name string) (*VirtualNetworkAdapter, error) {
	wql := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s'", name)
	ins, err := session.FindFirstInstance(wql)
	if err != nil {
		return nil, err
	}
	return NewVirtualNetworkAdapterFromInstance(ins)
}

func (vna *VirtualNetworkAdapter) GetGuestNetworkAdapterConfiguration() (*GuestNetworkAdapterConfiguration, error) {
	inst, err := vna.GetRelated(Msvm_GuestNetworkAdapterConfiguration)
	if err != nil {
		return nil, err
	}
	return NewGuestNetworkAdapterConfiguration(inst)
}
