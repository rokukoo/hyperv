package switch_extension

import "github.com/rokukoo/hyperv/pkg/wmiext"

const (
	Msvm_EthernetSwitchPortVlanSettingData = "Msvm_EthernetSwitchPortVlanSettingData"
)

type EthernetSwitchPortVlanSettingData struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID           string   `json:"instance_id"`
	Caption              string   `json:"caption"`
	Description          string   `json:"description"`
	ElementName          string   `json:"element_name"`
	OperationMode        uint32   `json:"operation_mode"`
	AccessVlanId         uint16   `json:"access_vlan_id"`
	NativeVlanId         uint16   `json:"native_vlan_id"`
	PvlanMode            uint32   `json:"pvlan_mode"`
	PrimaryVlanId        uint16   `json:"primary_vlan_id"`
	SecondaryVlanId      uint16   `json:"secondary_vlan_id"`
	PruneVlanIdArray     []uint16 `json:"prune_vlan_id_array,omitempty"`
	TrunkVlanIdArray     []uint16 `json:"trunk_vlan_id_array,omitempty"`
	SecondaryVlanIdArray []uint16 `json:"secondary_vlan_id_array,omitempty"`

	*wmiext.Instance
}

func (espvsd *EthernetSwitchPortVlanSettingData) Path() string {
	return espvsd.S__PATH
}

func NewEthernetSwitchPortVlanSettingData(instance *wmiext.Instance) (*EthernetSwitchPortVlanSettingData, error) {
	espvsd := &EthernetSwitchPortVlanSettingData{}
	if err := instance.GetAll(espvsd); err != nil {
		return nil, err
	}
	return espvsd, nil
}
