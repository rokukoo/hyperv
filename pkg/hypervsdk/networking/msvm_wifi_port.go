package networking

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	Msvm_WiFiPort = "Msvm_WiFiPort"
)

// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-wifiport
type WiFiPort struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID                       string
	Caption                          string
	Description                      string
	ElementName                      string
	InstallDate                      string
	Name                             string
	OperationalStatus                []uint16
	StatusDescriptions               []string
	Status                           string
	HealthState                      uint16
	CommunicationStatus              uint16
	DetailedStatus                   uint16
	OperatingStatus                  uint16
	PrimaryStatus                    uint16
	EnabledState                     uint16
	OtherEnabledState                string
	RequestedState                   uint16
	EnabledDefault                   uint16
	TimeOfLastStateChange            string
	AvailableRequestedStates         []uint16
	TransitioningToState             uint16
	SystemCreationClassName          string
	SystemName                       string
	CreationClassName                string
	DeviceID                         string
	PowerManagementSupported         bool
	PowerManagementCapabilities      []uint16
	Availability                     uint16
	StatusInfo                       uint16
	LastErrorCode                    uint32
	ErrorDescription                 string
	ErrorCleared                     bool
	OtherIdentifyingInfo             []string
	PowerOnHours                     uint64
	TotalPowerOnHours                uint64
	IdentifyingDescriptions          []string
	AdditionalAvailability           []uint16
	MaxQuiesceTime                   uint64
	Speed                            uint64
	MaxSpeed                         uint64
	RequestedSpeed                   uint64
	UsageRestriction                 uint16
	PortType                         uint16
	OtherPortType                    string
	OtherNetworkPortType             string
	PortNumber                       uint16
	LinkTechnology                   uint16
	OtherLinkTechnology              string
	PermanentAddress                 string
	NetworkAddresses                 []string
	FullDuplex                       bool
	AutoSense                        bool
	SupportedMaximumTransmissionUnit uint64
	ActiveMaximumTransmissionUnit    uint64
	IsBound                          bool

	*wmiext.Instance
}

func (wp *WiFiPort) Path() string {
	return wp.S__PATH
}

func GetWiFiPort(con *wmiext.Service, ethernetName string) (*WiFiPort, error) {
	wifiPort := &WiFiPort{}
	wquery := fmt.Sprintf("SELECT * FROM Msvm_WiFiPort WHERE ElementName = '%s'", ethernetName)
	return wifiPort, con.FindFirstObject(wquery, wifiPort)
}
