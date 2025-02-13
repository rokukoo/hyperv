package hypervsdk

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

// https://learn.microsoft.com/zh-cn/windows/win32/fwp/wmi/netadaptercimprov/msft-netadapter
type NetworkAdapter struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	Caption                                          string
	Description                                      string
	InstallDate                                      string
	Name                                             string
	Status                                           string
	Availability                                     uint16
	ConfigManagerErrorCode                           uint32
	ConfigManagerUserConfig                          bool
	CreationClassName                                string
	DeviceID                                         string
	ErrorCleared                                     bool
	ErrorDescription                                 string
	LastErrorCode                                    uint32
	PNPDeviceID                                      string
	PowerManagementCapabilities                      []uint16
	PowerManagementSupported                         bool
	StatusInfo                                       uint16
	SystemCreationClassName                          string
	SystemName                                       string
	Speed                                            uint64
	MaxSpeed                                         uint64
	RequestedSpeed                                   uint64
	UsageRestriction                                 uint16
	PortType                                         uint16
	OtherPortType                                    string
	OtherNetworkPortType                             string
	PortNumber                                       uint16
	LinkTechnology                                   uint16
	OtherLinkTechnology                              string
	PermanentAddress                                 string
	NetworkAddresses                                 []string
	FullDuplex                                       bool
	AutoSense                                        bool
	SupportedMaximumTransmissionUnit                 uint64
	ActiveMaximumTransmissionUnit                    uint64
	InterfaceDescription                             string
	InterfaceName                                    string
	NetLuid                                          uint64
	InterfaceGuid                                    string
	InterfaceIndex                                   uint32
	DeviceName                                       string
	NetLuidIndex                                     uint32
	Virtual                                          bool
	Hidden                                           bool
	NotUserRemovable                                 bool
	IMFilter                                         bool
	InterfaceType                                    uint32
	HardwareInterface                                bool
	WdmInterface                                     bool
	EndPointInterface                                bool
	iSCSIInterface                                   bool
	State                                            uint32
	NdisMedium                                       uint32
	NdisPhysicalMedium                               uint32
	InterfaceOperationalStatus                       uint32
	OperationalStatusDownDefaultPortNotAuthenticated bool
	OperationalStatusDownMediaDisconnected           bool
	OperationalStatusDownInterfacePaused             bool
	OperationalStatusDownLowPowerState               bool
	InterfaceAdminStatus                             uint32
	MediaConnectState                                uint32
	MtuSize                                          uint32
	VlanID                                           uint16
	TransmitLinkSpeed                                uint64
	ReceiveLinkSpeed                                 uint64
	PromiscuousMode                                  bool
	DeviceWakeUpEnable                               bool
	ConnectorPresent                                 bool
	MediaDuplexState                                 uint32
	DriverDate                                       string
	DriverDateData                                   uint64
	DriverVersionString                              string
	DriverName                                       string
	DriverDescription                                string
	MajorDriverVersion                               uint16
	MinorDriverVersion                               uint16
	DriverMajorNdisVersion                           uint8
	DriverMinorNdisVersion                           uint8
	PnPDeviceID                                      string
	DriverProvider                                   string
	ComponentID                                      string
	LowerLayerInterfaceIndices                       []uint32
	HigherLayerInterfaceIndices                      []uint32
	AdminLocked                                      bool

	*wmiext.Instance
}

func GetNetworkAdapterByName(name string) (adapter *NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StadardCimV2)
	if err != nil {
		return nil, err
	}
	netadapter := &NetworkAdapter{}
	wquery := fmt.Sprintf("SELECT * FROM MSFT_NetAdapter WHERE Name = '%s'", name)
	return netadapter, con.FindFirstObject(wquery, netadapter)
}

func GetNetworkAdapterByDescription(description string) (adapter *NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StadardCimV2)
	if err != nil {
		return nil, err
	}
	netadapter := &NetworkAdapter{}
	wquery := fmt.Sprintf("SELECT * FROM MSFT_NetAdapter WHERE Description = '%s'", description)
	return netadapter, con.FindFirstObject(wquery, netadapter)
}
