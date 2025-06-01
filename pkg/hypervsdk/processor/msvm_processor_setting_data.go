package processor

import (
	"fmt"

	"github.com/rokukoo/hyperv/pkg/hypervsdk"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

const (
	Msvm_ProcessorSettingData = "Msvm_ProcessorSettingData"
)

type ProcessorSettingData struct {
	S__PATH                        string
	InstanceID                     string
	Caption                        string // = "Processor"
	Description                    string // = "A logical processor of the hypervisor running on the host computer system."
	ElementName                    string
	ResourceType                   uint16 // = 3
	OtherResourceType              string
	ResourceSubType                string // = "Microsoft:Hyper-V:Processor"
	PoolID                         string
	ConsumerVisibility             uint16
	HostResource                   []string
	AllocationUnits                string // = "percent / 1000"
	VirtualQuantity                uint64 // = "count"
	Reservation                    uint64 // = 0
	Limit                          uint64 // = 100000
	Weight                         uint32 // = 100
	AutomaticAllocation            bool   // = True
	AutomaticDeallocation          bool   // = True
	Parent                         string
	Connection                     []string
	Address                        string
	MappingBehavior                uint16
	AddressOnParent                string
	VirtualQuantityUnits           string // = "count"
	LimitCPUID                     bool
	HwThreadsPerCore               uint64
	LimitProcessorFeatures         bool
	MaxProcessorsPerNumaNode       uint64
	MaxNumaNodesPerSocket          uint64
	EnableHostResourceProtection   bool
	CpuGroupId                     string
	HideHypervisorPresent          bool
	ExposeVirtualizationExtensions bool

	//service          *virtualsystem.VirtualSystemManagementService `json:"-"`
	*wmiext.Instance `json:"-"`
}

const ProcessorResourceType = "Microsoft:Hyper-V:Processor"

func (psd *ProcessorSettingData) Path() string {
	return psd.S__PATH
}

func GetDefaultProcessorSettingData() (*ProcessorSettingData, error) {
	settings := &ProcessorSettingData{}
	return settings, hypervsdk.PopulateDefaults(ProcessorResourceType, settings)
}

func CreateProcessorSettings(settings *ProcessorSettingData) (string, error) {
	str, err := hypervsdk.CreateResourceSettingGeneric(settings, ProcessorResourceType)
	if err != nil {
		err = fmt.Errorf("could not create processor settings: %w", err)
	}
	return str, err
}
