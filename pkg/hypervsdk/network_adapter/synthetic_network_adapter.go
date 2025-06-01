package network_adapter

import (
	"github.com/rokukoo/hyperv/pkg/hypervsdk/networking"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

type SyntheticNetworkAdapter struct {
	*networking.SyntheticEthernetPortSettingData
}

// SyntheticNetworkAdapterFromInstance creates a new SyntheticNetworkAdapter instance from a WMI instance
func NewSyntheticNetworkAdapterFromInstance(instance *wmiext.Instance) (*SyntheticNetworkAdapter, error) {
	syntheticEthernetPortSettingData := &networking.SyntheticEthernetPortSettingData{}
	if err := instance.GetAll(syntheticEthernetPortSettingData); err != nil {
		return nil, err
	}
	return &SyntheticNetworkAdapter{syntheticEthernetPortSettingData}, nil
}

func (sna *SyntheticNetworkAdapter) SetElementName(elementName string) error {
	sna.ElementName = elementName
	return sna.Put("ElementName", elementName)
}

func (sna *SyntheticNetworkAdapter) GetEthernetPortAllocationSettingData() (epas *networking.EthernetPortAllocationSettingData, err error) {
	inst, err := sna.GetRelated("Msvm_EthernetPortAllocationSettingData")
	if err != nil {
		return
	}
	return networking.NewEthernetPortAllocationSettingDataFromInstance(inst)
}
