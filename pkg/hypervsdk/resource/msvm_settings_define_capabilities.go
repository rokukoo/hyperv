package resource

import "github.com/rokukoo/hyperv/pkg/wmiext"

const (
	Msvm_SettingsDefineCapabilities = "Msvm_SettingsDefineCapabilities"
)

// SettingsDefineCapabilities represents the settings define capabilities in Go.
type SettingsDefineCapabilities struct {
	SupportStatement uint16 `json:"support_statement"`
	GroupComponent   string `json:"group_component"`
	PartComponent    string `json:"part_component"`
	PropertyPolicy   uint16 `json:"property_policy"`
	ValueRole        uint16 `json:"value_role"`
	ValueRange       uint16 `json:"value_range"`

	// Instance is the WMI instance.
	*wmiext.Instance `json:"-"`
}
