package switch_extension

import "github.com/rokukoo/hypervctl/pkg/wmiext"

// EthernetSwitchFeatureCapabilities 代表 Hyper-V 以太网交换机功能能力
type EthernetSwitchFeatureCapabilities struct {
	S__PATH string `json:"-"`

	InstanceID    string `json:"instance_id"`
	Caption       string `json:"caption" default:"Ethernet Switch Feature Capabilities"`
	Description   string `json:"description" default:"Microsoft Virtual Ethernet Switch Feature Capabilities"`
	ElementName   string `json:"element_name" default:"Ethernet Switch Port Bandwidth Settings"`
	FeatureID     string `json:"feature_id"`
	Applicability uint16 `json:"applicability"`
	Version       string `json:"version"`

	*wmiext.Instance `json:"-"`
}

func (esc *EthernetSwitchFeatureCapabilities) Path() string {
	return esc.S__PATH
}

func NewEthernetSwitchFeatureCapabilities(inst *wmiext.Instance) (*EthernetSwitchFeatureCapabilities, error) {
	esc := &EthernetSwitchFeatureCapabilities{}
	return esc, inst.GetAll(esc)
}
