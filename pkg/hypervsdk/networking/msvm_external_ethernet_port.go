package networking

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	Msvm_ExternalEthernetPort = "Msvm_ExternalEthernetPort"
)

type ExternalEthernetPort struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID                       string   `json:"instance_id"`
	Caption                          string   `json:"caption"`
	Description                      string   `json:"description"`
	ElementName                      string   `json:"element_name"`
	InstallDate                      string   `json:"install_date,omitempty"` // datetime 采用 string
	Name                             string   `json:"name"`
	OperationalStatus                []uint16 `json:"operational_status"`
	StatusDescriptions               []string `json:"status_descriptions"`
	Status                           string   `json:"status"`
	HealthState                      uint16   `json:"health_state"`
	CommunicationStatus              uint16   `json:"communication_status,omitempty"`
	DetailedStatus                   uint16   `json:"detailed_status,omitempty"`
	OperatingStatus                  uint16   `json:"operating_status,omitempty"`
	PrimaryStatus                    uint16   `json:"primary_status,omitempty"`
	EnabledState                     uint16   `json:"enabled_state"`
	OtherEnabledState                string   `json:"other_enabled_state,omitempty"`
	RequestedState                   uint16   `json:"requested_state"`
	EnabledDefault                   uint16   `json:"enabled_default"`
	TimeOfLastStateChange            string   `json:"time_of_last_state_change,omitempty"`
	AvailableRequestedStates         []uint16 `json:"available_requested_states,omitempty"`
	TransitioningToState             uint16   `json:"transitioning_to_state,omitempty"`
	SystemCreationClassName          string   `json:"system_creation_class_name"`
	SystemName                       string   `json:"system_name"`
	CreationClassName                string   `json:"creation_class_name"`
	DeviceID                         string   `json:"device_id"`
	PowerManagementSupported         bool     `json:"power_management_supported"`
	PowerManagementCapabilities      []uint16 `json:"power_management_capabilities,omitempty"`
	Availability                     uint16   `json:"availability,omitempty"`
	StatusInfo                       uint16   `json:"status_info,omitempty"`
	LastErrorCode                    uint32   `json:"last_error_code,omitempty"`
	ErrorDescription                 string   `json:"error_description,omitempty"`
	ErrorCleared                     bool     `json:"error_cleared"`
	OtherIdentifyingInfo             []string `json:"other_identifying_info,omitempty"`
	PowerOnHours                     uint64   `json:"power_on_hours,omitempty"`
	TotalPowerOnHours                uint64   `json:"total_power_on_hours,omitempty"`
	IdentifyingDescriptions          []string `json:"identifying_descriptions,omitempty"`
	AdditionalAvailability           []string `json:"additional_availability,omitempty"`
	MaxQuiesceTime                   uint64   `json:"max_quiesce_time,omitempty"`
	Speed                            uint64   `json:"speed,omitempty"`
	MaxSpeed                         uint64   `json:"max_speed,omitempty"`
	RequestedSpeed                   uint64   `json:"requested_speed,omitempty"`
	UsageRestriction                 uint16   `json:"usage_restriction,omitempty"`
	PortType                         uint16   `json:"port_type,omitempty"`
	OtherPortType                    string   `json:"other_port_type,omitempty"`
	OtherNetworkPortType             string   `json:"other_network_port_type,omitempty"`
	PortNumber                       uint16   `json:"port_number,omitempty"`
	LinkTechnology                   uint16   `json:"link_technology,omitempty"`
	OtherLinkTechnology              string   `json:"other_link_technology,omitempty"`
	PermanentAddress                 string   `json:"permanent_address,omitempty"`
	NetworkAddresses                 []string `json:"network_addresses,omitempty"`
	FullDuplex                       bool     `json:"full_duplex"`
	AutoSense                        bool     `json:"auto_sense"`
	SupportedMaximumTransmissionUnit uint64   `json:"supported_maximum_transmission_unit,omitempty"`
	ActiveMaximumTransmissionUnit    uint64   `json:"active_maximum_transmission_unit,omitempty"`
	MaxDataSize                      uint32   `json:"max_data_size,omitempty"`
	Capabilities                     []uint16 `json:"capabilities,omitempty"`
	CapabilityDescriptions           []string `json:"capability_descriptions,omitempty"`
	EnabledCapabilities              []uint16 `json:"enabled_capabilities,omitempty"`
	OtherEnabledCapabilities         []string `json:"other_enabled_capabilities,omitempty"`
	IsBound                          bool     `json:"is_bound"`

	*wmiext.Instance
}

func (eep *ExternalEthernetPort) Path() string {
	return eep.S__PATH
}

func GetExternalEthernetPort(con *wmiext.Service, ethernetName string) (*ExternalEthernetPort, error) {
	extPort := &ExternalEthernetPort{}
	wquery := fmt.Sprintf("SELECT * FROM Msvm_ExternalEthernetPort WHERE ElementName = '%s'", ethernetName)
	return extPort, con.FindFirstObject(wquery, extPort)
}

func ListEnabledExternalEthernetPort(session *wmiext.Service) ([]*ExternalEthernetPort, error) {
	var extPorts []*ExternalEthernetPort
	wquery := fmt.Sprintf("SELECT * FROM Msvm_ExternalEthernetPort WHERE EnabledState = 2")
	return extPorts, session.FindObjects(wquery, extPorts)
}
