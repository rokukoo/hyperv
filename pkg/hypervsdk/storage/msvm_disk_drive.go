package storage

import (
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

const (
	Msvm_DiskDrive = "Msvm_DiskDrive"
)

type DiskDrive struct {
	InstanceID                  string    `json:"instance_id"`
	Caption                     string    `json:"caption"`
	Description                 string    `json:"description"`
	ElementName                 string    `json:"element_name"`
	InstallDate                 time.Time `json:"install_date"`
	Name                        string    `json:"name"`
	OperationalStatus           []uint16  `json:"operational_status"`
	StatusDescriptions          []string  `json:"status_descriptions"`
	Status                      string    `json:"status"`
	HealthState                 uint16    `json:"health_state"`
	CommunicationStatus         uint16    `json:"communication_status"`
	DetailedStatus              uint16    `json:"detailed_status"`
	OperatingStatus             uint16    `json:"operating_status"`
	PrimaryStatus               uint16    `json:"primary_status"`
	EnabledState                uint16    `json:"enabled_state"`
	OtherEnabledState           string    `json:"other_enabled_state"`
	RequestedState              uint16    `json:"requested_state"`
	EnabledDefault              uint16    `json:"enabled_default"`
	TimeOfLastStateChange       time.Time `json:"time_of_last_state_change"`
	AvailableRequestedStates    []uint16  `json:"available_requested_states"`
	TransitioningToState        uint16    `json:"transitioning_to_state"`
	SystemCreationClassName     string    `json:"system_creation_class_name"`
	SystemName                  string    `json:"system_name"`
	CreationClassName           uint16    `json:"creation_class_name"`
	DeviceID                    string    `json:"device_id"`
	PowerManagementSupported    bool      `json:"power_management_supported"`
	PowerManagementCapabilities []uint16  `json:"power_management_capabilities"`
	Availability                uint16    `json:"availability"`
	StatusInfo                  uint16    `json:"status_info"`
	LastErrorCode               uint32    `json:"last_error_code"`
	ErrorDescription            string    `json:"error_description"`
	ErrorCleared                bool      `json:"error_cleared"`
	OtherIdentifyingInfo        []string  `json:"other_identifying_info"`
	PowerOnHours                uint64    `json:"power_on_hours"`
	TotalPowerOnHours           uint64    `json:"total_power_on_hours"`
	IdentifyingDescriptions     []string  `json:"identifying_descriptions"`
	AdditionalAvailability      []uint16  `json:"additional_availability"`
	MaxQuiesceTime              uint64    `json:"max_quiesce_time"`
	Capabilities                []uint16  `json:"capabilities"`
	CapabilityDescriptions      []string  `json:"capability_descriptions"`
	ErrorMethodology            string    `json:"error_methodology"`
	CompressionMethod           string    `json:"compression_method"`
	NumberOfMediaSupported      uint32    `json:"number_of_media_supported"`
	MaxMediaSize                uint64    `json:"max_media_size"`
	DefaultBlockSize            uint64    `json:"default_block_size"`
	MaxBlockSize                uint64    `json:"max_block_size"`
	MinBlockSize                uint64    `json:"min_block_size"`
	NeedsCleaning               bool      `json:"needs_cleaning"`
	MediaIsLocked               bool      `json:"media_is_locked"`
	Security                    uint16    `json:"security"`
	LastCleaned                 time.Time `json:"last_cleaned"`
	MaxAccessTime               uint64    `json:"max_access_time"`
	UncompressedDataRate        uint32    `json:"uncompressed_data_rate"`
	LoadTime                    uint64    `json:"load_time"`
	UnloadTime                  uint64    `json:"unload_time"`
	MountCount                  uint64    `json:"mount_count"`
	TimeOfLastMount             time.Time `json:"time_of_last_mount"`
	TotalMountTime              uint64    `json:"total_mount_time"`
	UnitsDescription            string    `json:"units_description"`
	MaxUnitsBeforeCleaning      uint64    `json:"max_units_before_cleaning"`
	UnitsUsed                   uint64    `json:"units_used"`
	DriveNumber                 uint32    `json:"drive_number"`

	*wmiext.Instance `json:"-"`
}
