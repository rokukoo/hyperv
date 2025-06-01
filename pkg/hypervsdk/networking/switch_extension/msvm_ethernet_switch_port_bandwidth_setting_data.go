package switch_extension

import (
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

const (
	Msvm_EthernetSwitchPortBandwidthSettingData = "Msvm_EthernetSwitchPortBandwidthSettingData"
)

// EthernetSwitchPortBandwidthSettingData 代表以太网交换机端口带宽设置
type EthernetSwitchPortBandwidthSettingData struct {
	S__PATH string `json:"-"`

	InstanceID  string `json:"instance_id"`
	Caption     string `json:"caption" default:"Ethernet Switch Port Bandwidth Settings"`
	Description string `json:"description" default:"Represents the port bandwidth settings."`
	ElementName string `json:"element_name" default:"Ethernet Switch Port Bandwidth Settings"`
	Reservation uint64 `json:"reservation" default:"0"`
	Weight      uint64 `json:"weight" default:"0"`
	Limit       uint64 `json:"limit" default:"0"`
	BurstLimit  uint64 `json:"burst_limit" default:"0"`
	BurstSize   uint64 `json:"burst_size" default:"0"`

	*wmiext.Instance `json:"-"`
}

func (espbsd *EthernetSwitchPortBandwidthSettingData) Path() string {
	return espbsd.S__PATH
}

func (espbsd *EthernetSwitchPortBandwidthSettingData) SetLimit(limit uint64) error {
	espbsd.Limit = limit
	return espbsd.Put("Limit", limit)
}

func (espbsd *EthernetSwitchPortBandwidthSettingData) SetReservation(reservation uint64) error {
	espbsd.Reservation = reservation
	return espbsd.Put("Reservation", reservation)
}

func (espbsd *EthernetSwitchPortBandwidthSettingData) SetBurstLimit(burstLimit uint64) error {
	espbsd.BurstLimit = burstLimit
	return espbsd.Put("BurstLimit", burstLimit)
}

func (espbsd *EthernetSwitchPortBandwidthSettingData) SetBurstSize(burstSize uint64) error {
	espbsd.BurstSize = burstSize
	return espbsd.Put("BurstSize", burstSize)
}

func NewEthernetSwitchPortBandwidthSettingData(inst *wmiext.Instance) (*EthernetSwitchPortBandwidthSettingData, error) {
	espbsd := &EthernetSwitchPortBandwidthSettingData{}
	return espbsd, inst.GetAll(espbsd)
}
