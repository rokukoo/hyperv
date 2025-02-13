package networking

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

// VirtualEthernetSwitch represents a virtual ethernet switch
// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-virtualethernetswitch
type VirtualEthernetSwitch struct {
	S__PATH                     string `json:"-"`
	S__CLASS                    string `json:"-"`
	Caption                     string
	Description                 string
	ElementName                 string
	InstallDate                 time.Time
	OperationalStatus           []uint16
	Status                      string
	HealthState                 uint16
	EnabledState                uint16
	OtherEnabledState           string
	RequestedState              uint16
	TimeOfLastStateChange       time.Time
	Name                        string
	PrimaryOwnerName            string
	IdentifyingDescriptions     []string
	OtherIdentifyingInfo        []string
	Dedicated                   []uint16
	ResetCapability             uint16
	PowerManagementCapabilities []uint16
	StatusDescriptions          []string
	EnabledDefault              uint16
	CreationClassName           string
	PrimaryOwnerContact         string
	Roles                       []string
	NameFormat                  string
	OtherDedicatedDescriptions  []string
	ScopeOfResidence            string
	NumLearnableAddresses       uint32
	MaxVMQOffloads              uint32
	MaxChimneyOffloads          uint32

	*wmiext.Instance
}

func (vswitch *VirtualEthernetSwitch) Path() string {
	return vswitch.S__PATH
}

func (vswitch *VirtualEthernetSwitch) GetInternalPortAllocSettings() (*EthernetPortAllocationSettingData, error) {
	ethernetPortAllocSettings, err := vswitch.GetEthernetPortAllocSettings()
	if err != nil {
		return nil, err
	}
	for _, setting := range ethernetPortAllocSettings {
		hostResPath := setting.HostResource[0]
		object, err := vswitch.GetService().GetObject(hostResPath)
		if err != nil {
			return nil, err
		}
		className, err := object.GetClassName()
		if err != nil {
			return nil, err
		}
		if className == "Msvm_ComputerSystem" {
			return &setting, nil
		}
	}
	return nil, wmiext.NotFound
}

func (vswitch *VirtualEthernetSwitch) GetExternalPortAllocSettings() (*EthernetPortAllocationSettingData, error) {
	ethernetPortAllocSettings, err := vswitch.GetEthernetPortAllocSettings()
	if err != nil {
		return nil, err
	}

	for _, setting := range ethernetPortAllocSettings {
		hostResPath := setting.HostResource[0]
		object, err := vswitch.GetService().GetObject(hostResPath)
		if err != nil {
			return nil, err
		}
		className, err := object.GetClassName()
		if err != nil {
			return nil, err
		}
		if className == Msvm_ExternalEthernetPort || className == Msvm_WiFiPort {
			return &setting, nil
		}
	}
	return nil, wmiext.NotFound
}

func (vswitch *VirtualEthernetSwitch) GetEthernetPortAllocSettings() ([]EthernetPortAllocationSettingData, error) {
	var (
		settingData []EthernetPortAllocationSettingData
	)
	virtualEthernetSwitchSettingData, err := vswitch.ActiveVirtualEthernetSwitchSettingData()
	if err != nil {
		return nil, err
	}
	return settingData, vswitch.GetService().FindRelatedObjects(virtualEthernetSwitchSettingData.Path(), Msvm_EthernetPortAllocationSettingData, &settingData)
}

func (vswitch *VirtualEthernetSwitch) ActiveVirtualEthernetSwitchSettingData() (*VirtualEthernetSwitchSettingData, error) {
	port := &VirtualEthernetSwitchSettingData{}
	return port, vswitch.GetService().FindFirstRelatedObject(vswitch.Path(), "Msvm_VirtualEthernetSwitchSettingData", port)
}

func FirstVirtualEthernetSwitchByName(session *wmiext.Service, name string) (*VirtualEthernetSwitch, error) {
	vswitch := &VirtualEthernetSwitch{}
	wql := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE ElementName = '%s'", name)
	return vswitch, session.FindFirstObject(wql, vswitch)
}
