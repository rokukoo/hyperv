package win32

import (
	"fmt"
	"github.com/pkg/errors"
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

func (na *NetworkAdapter) Configure(
	ipAddress []string,
	subnetMask []string,
	gateway []string,
	dnsServer []string,
) (err error) {
	con, err := wmiext.NewLocalService(wmiext.CimV2)
	if err != nil {
		return
	}
	wquery := fmt.Sprintf("SELECT * FROM Win32_NetworkAdapterConfiguration WHERE InterfaceIndex = %d", na.InterfaceIndex)
	netAdapterConfiguration, err := con.FindFirstInstance(wquery)
	if err != nil {
		return
	}

	var returnValue int32
	if err = netAdapterConfiguration.Method("EnableStatic").
		In("IPAddress", ipAddress).
		In("SubnetMask", subnetMask).
		Execute().
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return
	}

	if returnValue == -2147024891 {
		return wmiext.PermissionDenied
	}

	if returnValue != 0 {
		return fmt.Errorf("failed to enable static ip, return value: %d", returnValue)
	}

	// 设置 网关
	gatewayCostMetrics := []uint16{1} // 网关跃点数
	if err = netAdapterConfiguration.Method("SetGateways").
		In("DefaultIPGateway", gateway).
		In("GatewayCostMetric", gatewayCostMetrics).
		Execute().
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return fmt.Errorf("failed to set gateway, return value: %d", returnValue)
	}

	if returnValue != 0 {
		return fmt.Errorf("failed to set gateway, return value: %d", returnValue)
	}

	// 设置 DNS
	if err = netAdapterConfiguration.Method("SetDNSServerSearchOrder").
		In("DNSServerSearchOrder", dnsServer).
		Execute().
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return fmt.Errorf("failed to set dnsServer, return value: %d", returnValue)
	}

	if returnValue != 0 {
		return errors.New("failed to set dnsServer")
	}

	return nil
}

func ListPhysicalNetAdapter() (adapters []NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StandardCimV2)
	if err != nil {
		return nil, err
	}
	wquery := "SELECT * FROM MSFT_NetAdapter WHERE Virtual = false"
	if err = con.FindObjects(wquery, &adapters); err != nil {
		return nil, err
	}
	return adapters, nil
}

func ListNetworkAdapters() (adapters []NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StandardCimV2)
	if err != nil {
		return nil, err
	}
	wquery := "SELECT * FROM MSFT_NetAdapter"
	if err = con.FindObjects(wquery, &adapters); err != nil {
		return nil, err
	}
	return adapters, nil
}

func GetNetworkAdapterByName(name string) (adapter *NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StandardCimV2)
	if err != nil {
		return nil, err
	}
	netadapter := &NetworkAdapter{}
	wquery := fmt.Sprintf("SELECT * FROM MSFT_NetAdapter WHERE Name = '%s'", name)
	return netadapter, con.FindFirstObject(wquery, netadapter)
}

func GetNetworkAdapterByDescription(description string) (adapter *NetworkAdapter, err error) {
	con, err := wmiext.NewLocalService(wmiext.StandardCimV2)
	if err != nil {
		return nil, err
	}
	netadapter := &NetworkAdapter{}
	wquery := fmt.Sprintf("SELECT * FROM MSFT_NetAdapter WHERE InterfaceDescription = '%s'", description)
	return netadapter, con.FindFirstObject(wquery, netadapter)
}
