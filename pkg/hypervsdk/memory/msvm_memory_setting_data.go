package memory

import (
	"fmt"

	"github.com/rokukoo/hyperv/pkg/hypervsdk"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

const (
	MemoryResourceType     = "Microsoft:Hyper-V:Memory"
	Msvm_MemorySettingData = "Msvm_MemorySettingData"
)

type MemorySettingsData struct {
	S__PATH                    string
	InstanceID                 string
	Caption                    string // = "Memory Default Settings"
	Description                string // = "Describes the default settings for the memory resources."
	ElementName                string
	ResourceType               uint16 // = 4
	OtherResourceType          string
	ResourceSubType            string // = "Microsoft:Hyper-V:Memory"
	PoolID                     string
	ConsumerVisibility         uint16
	HostResource               []string
	HugePagesEnabled           bool
	AllocationUnits            string // = "byte * 2^20"
	VirtualQuantity            uint64
	Reservation                uint64
	Limit                      uint64
	Weight                     uint32
	AutomaticAllocation        bool // = True
	AutomaticDeallocation      bool // = True
	Parent                     string
	Connection                 []string
	Address                    string
	MappingBehavior            uint16
	AddressOnParent            string
	VirtualQuantityUnits       string // = "byte * 2^20"
	DynamicMemoryEnabled       bool
	TargetMemoryBuffer         uint32
	IsVirtualized              bool // = True
	SwapFilesInUse             bool
	MaxMemoryBlocksPerNumaNode uint64
	SgxSize                    uint64
	SgxEnabled                 bool

	//service *virtualsystem.VirtualSystemManagementService `json:"-"`
	*wmiext.Instance
}

func (msd *MemorySettingsData) Path() string {
	return msd.S__PATH
}

func GetDefaultMemorySettingsData() (*MemorySettingsData, error) {
	settings := &MemorySettingsData{}
	return settings, hypervsdk.PopulateDefaults(MemoryResourceType, settings)
}

func CreateMemorySettings(settings *MemorySettingsData) (string, error) {
	str, err := hypervsdk.CreateResourceSettingGeneric(settings, MemoryResourceType)
	if err != nil {
		err = fmt.Errorf("could not create memory settings: %w", err)
	}
	return str, err
}
