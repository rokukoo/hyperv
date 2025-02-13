package networking

import (
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

const (
	Msvm_VirtualEthernetSwitchSettingData = "Msvm_VirtualEthernetSwitchSettingData"
)

type VirtualEthernetSwitchSettingData struct {
	S__PATH                              string `json:"-"`
	S__CLASS                             string `json:"-"`
	InstanceID                           string
	Caption                              string
	Description                          string
	ElementName                          string
	VirtualSystemIdentifier              string
	VirtualSystemType                    string
	Notes                                []string
	CreationTime                         time.Time
	ConfigurationID                      string
	ConfigurationDataRoot                string
	ConfigurationFile                    string
	SnapshotDataRoot                     string
	SuspendDataRoot                      string
	SwapFileDataRoot                     string
	LogDataRoot                          string
	AutomaticStartupAction               uint16
	AutomaticStartupActionDelay          time.Duration
	AutomaticStartupActionSequenceNumber uint16
	AutomaticShutdownAction              uint16
	AutomaticRecoveryAction              uint16
	RecoveryFile                         string
	VLANConnection                       []string
	AssociatedResourcePool               []string
	MaxNumMACAddress                     uint32
	IOVPreferred                         bool
	ExtensionOrder                       []string
	BandwidthReservationMode             uint32
	TeamingEnabled                       bool
	PacketDirectEnabled                  bool

	*wmiext.Instance
}

func (vesd *VirtualEthernetSwitchSettingData) Path() string {
	return vesd.S__PATH
}
