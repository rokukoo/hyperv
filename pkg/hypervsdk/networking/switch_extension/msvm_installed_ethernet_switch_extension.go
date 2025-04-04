package switch_extension

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

const (
	Msvm_EthernetSwitchFeatureCapabilities = "Msvm_EthernetSwitchFeatureCapabilities"
	Msvm_InstalledEthernetSwitchExtension  = "Msvm_InstalledEthernetSwitchExtension"
)

// InstalledEthernetSwitchExtension 代表 Hyper-V 安装的以太网交换机扩展
type InstalledEthernetSwitchExtension struct {
	S__PATH string `json:"-"`

	InstanceID          string    `json:"instance_id"`
	Caption             string    `json:"caption" default:"System Virtual Ethernet Switch Extension"`
	Description         string    `json:"description" default:"Microsoft NDIS Packet Capture Filter Driver"`
	ElementName         string    `json:"element_name" default:"Microsoft NDIS Capture"`
	InstallDate         time.Time `json:"install_date"`
	Name                string    `json:"name"`
	OperationalStatus   []uint16  `json:"operational_status" default:"[2]"`
	StatusDescriptions  []string  `json:"status_descriptions" default:"[OK]"`
	Status              string    `json:"status" default:"OK"`
	HealthState         uint16    `json:"health_state" default:"5"`
	CommunicationStatus uint16    `json:"communication_status"`
	DetailedStatus      uint16    `json:"detailed_status"`
	OperatingStatus     uint16    `json:"operating_status"`
	PrimaryStatus       uint16    `json:"primary_status"`
	ExtensionType       uint8     `json:"extension_type"`
	Vendor              string    `json:"vendor"`
	Version             string    `json:"version"`

	*wmiext.Instance `json:"-"`
}

func (iese *InstalledEthernetSwitchExtension) Path() string {
	return iese.S__PATH
}

// NewInstalledEthernetSwitchExtension 创建一个新的 InstalledEthernetSwitchExtension 实例
func NewInstalledEthernetSwitchExtension(inst *wmiext.Instance) (*InstalledEthernetSwitchExtension, error) {
	iese := &InstalledEthernetSwitchExtension{}
	return iese, inst.GetAll(iese)
}

func (iese *InstalledEthernetSwitchExtension) GetFeatureCapabilities() ([]*wmiext.Instance, error) {
	return iese.GetAllRelated(Msvm_EthernetSwitchFeatureCapabilities)
}

func (iese *InstalledEthernetSwitchExtension) GetFeatureCapabilityByName(name string) (*EthernetSwitchFeatureCapabilities, error) {
	caps, err := iese.GetFeatureCapabilities()
	if err != nil {
		return nil, err
	}

	for _, capability := range caps {
		capName, err := capability.GetAsString("ElementName")
		if err != nil {
			return nil, err
		}
		if capName == name {
			foundInstance, err := capability.CloneInstance()
			if err != nil {
				return nil, err
			}
			return NewEthernetSwitchFeatureCapabilities(foundInstance)
		}
	}

	return nil, errors.Wrapf(wmiext.NotFound, "Unable to find Feature Capability [%s]", name)
}
