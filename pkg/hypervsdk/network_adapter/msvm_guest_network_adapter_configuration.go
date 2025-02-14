package network_adapter

import "github.com/rokukoo/hypervctl/pkg/wmiext"

const (
	Msvm_GuestNetworkAdapterConfiguration = "Msvm_GuestNetworkAdapterConfiguration"
)

type ProtocolIFType = uint16

const (
	Unknown ProtocolIFType = iota
	Other
	IPv4   = 4096
	IPv6   = 4097
	IPv4v6 = 4098
)

// GuestNetworkAdapterConfiguration represents a network adapter configuration for a guest
//
// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-guestnetworkadapterconfiguration
type GuestNetworkAdapterConfiguration struct {
	S__PATH string `json:"-"`

	InstanceID       string   `json:"instance_id"`
	ProtocolIFType   uint16   `json:"protocol_if_type"`
	DHCPEnabled      bool     `json:"dhcp_enabled"`
	IPAddresses      []string `json:"ip_addresses"`
	Subnets          []string `json:"subnets"`
	DefaultGateways  []string `json:"default_gateways"`
	DNSServers       []string `json:"dns_servers"`
	IPAddressOrigins []uint16 `json:"ip_address_origins"`

	*wmiext.Instance `json:"-"`
}

func (gnac *GuestNetworkAdapterConfiguration) Path() string {
	return gnac.S__PATH
}

func NewGuestNetworkAdapterConfiguration(inst *wmiext.Instance) (*GuestNetworkAdapterConfiguration, error) {
	guestNetworkAdapterConfiguration := &GuestNetworkAdapterConfiguration{}
	return guestNetworkAdapterConfiguration, inst.GetAll(guestNetworkAdapterConfiguration)
}

func (gnac *GuestNetworkAdapterConfiguration) SetDHCPEnabled(dhcpEnabled bool) error {
	gnac.DHCPEnabled = dhcpEnabled
	return gnac.Put("DHCPEnabled", dhcpEnabled)
}

func (gnac *GuestNetworkAdapterConfiguration) SetProtocolIFType(protocolIFType ProtocolIFType) error {
	gnac.ProtocolIFType = protocolIFType
	return gnac.Put("ProtocolIFType", protocolIFType)
}

func (gnac *GuestNetworkAdapterConfiguration) SetIPAddresses(ipAddresses []string) error {
	gnac.IPAddresses = ipAddresses
	return gnac.Put("IPAddresses", ipAddresses)
}

func (gnac *GuestNetworkAdapterConfiguration) SetSubnets(subnets []string) error {
	gnac.Subnets = subnets
	return gnac.Put("Subnets", subnets)
}

func (gnac *GuestNetworkAdapterConfiguration) SetDefaultGateways(defaultGateways []string) error {
	gnac.DefaultGateways = defaultGateways
	return gnac.Put("DefaultGateways", defaultGateways)
}

func (gnac *GuestNetworkAdapterConfiguration) SetDNSServers(dnsServers []string) error {
	gnac.DNSServers = dnsServers
	return gnac.Put("DNSServers", dnsServers)
}
