package virtual_system

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/memory"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/network_adapter"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/processor"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/controller"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/disk"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/drive"
	utils "github.com/rokukoo/hypervctl/pkg/hypervsdk/utils"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"time"
)

const (
	Msvm_ComputerSystem = "Msvm_ComputerSystem"
)

type ComputerSystem struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	InstanceID                               string
	Caption                                  string
	Description                              string
	ElementName                              string
	InstallDate                              time.Time
	OperationalStatus                        []uint16
	StatusDescriptions                       []string
	Status                                   string
	HealthState                              uint16
	CommunicationStatus                      uint16
	DetailedStatus                           uint16
	OperatingStatus                          uint16
	PrimaryStatus                            uint16
	EnabledState                             uint16
	OtherEnabledState                        string
	RequestedState                           uint16
	EnabledDefault                           uint16
	TimeOfLastStateChange                    string
	AvailableRequestedStates                 []uint16
	TransitioningToState                     uint16
	CreationClassName                        string
	Name                                     string
	PrimaryOwnerName                         string
	PrimaryOwnerContact                      string
	Roles                                    []string
	NameFormat                               string
	OtherIdentifyingInfo                     []string
	IdentifyingDescriptions                  []string
	Dedicated                                []uint16
	OtherDedicatedDescriptions               []string
	ResetCapability                          uint16
	PowerManagementCapabilities              []uint16
	OnTimeInMilliseconds                     uint64
	ProcessID                                uint32
	TimeOfLastConfigurationChange            string
	NumberOfNumaNodes                        uint16
	ReplicationState                         uint16
	ReplicationHealth                        uint16
	ReplicationMode                          uint16
	FailedOverReplicationType                uint16
	LastReplicationType                      uint16
	LastApplicationConsistentReplicationTime string
	LastReplicationTime                      time.Time
	LastSuccessfulBackupTime                 string
	EnhancedSessionModeState                 uint16

	*wmiext.Instance `json:"-"`
}

func (vm *ComputerSystem) Path() string {
	return vm.S__PATH
}

func (vm *ComputerSystem) NewSCSIController() (*resource.ResourceAllocationSettingData, error) {
	rp, err := resource.GetPrimordialResourcePool(vm.GetService(), resource.ResourcePool_ResourceType_Parallel_SCSI_HBA)
	if err != nil {
		return nil, err
	}
	return rp.GetDefaultResourceAllocationSettingData()
}

func (vm *ComputerSystem) GetResourceAllocationSettingData(rtype resource.ResourcePool_ResourceType) ([]*resource.ResourceAllocationSettingData, error) {
	var (
		systemSettingData              *VirtualSystemSettingData
		resourceAllocationSettingDatas []*wmiext.Instance
		resourceAllocationSettingData  resource.ResourceAllocationSettingData

		resourceAllocationSettingDataCloneInst *wmiext.Instance

		curResType *resource.ResourceTypeValue
		err        error
		col        []*resource.ResourceAllocationSettingData
	)

	col = []*resource.ResourceAllocationSettingData{}

	if systemSettingData, err = vm.GetVirtualSystemSettingData(); err != nil {
		return nil, err
	}

	if resourceAllocationSettingDatas, err = systemSettingData.GetAllRelated(resource.Msvm_ResourceAllocationSettingData); err != nil {
		return nil, err
	}

	resType := resource.GetResourceTypeValue(rtype)
	for _, resourceAllocationSettingDataInst := range resourceAllocationSettingDatas {
		resourceAllocationSettingData = resource.ResourceAllocationSettingData{}
		if err = resourceAllocationSettingDataInst.GetAll(&resourceAllocationSettingData); err != nil {
			return nil, err
		}

		if curResType, err = resourceAllocationSettingData.GetResourceType(); err != nil {
			return nil, err
		}

		if !curResType.Equals(resType) {
			continue
		}

		if resourceAllocationSettingDataCloneInst, err = resourceAllocationSettingData.CloneInstance(); err != nil {
			return nil, err
		}

		resourceAllocationSettingData = resource.ResourceAllocationSettingData{}
		if err = resourceAllocationSettingDataCloneInst.GetAll(&resourceAllocationSettingData); err != nil {
			return nil, err
		}
		col = append(col, &resourceAllocationSettingData)
		if len(col) == 0 {
			return nil, errors.Wrapf(wmiext.NotFound, "GetResourceAllocationSettingData [%s] ", resType)
		}
		return col, nil
	}
	return nil, err
}

func (vm *ComputerSystem) GetSCSIControllers() (col []*resource.ResourceAllocationSettingData, err error) {
	col, err = vm.GetResourceAllocationSettingData(resource.ResourcePool_ResourceType_Parallel_SCSI_HBA)
	return
}

// GetProcessorSettingData returns the processor setting data of the Virtual Machine
func (vm *ComputerSystem) GetProcessorSettingData() (*processor.ProcessorSettingData, error) {
	setting, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	var processorSettingData = processor.ProcessorSettingData{
		//service: vm.service,
	}
	return &processorSettingData, vm.GetService().FindFirstRelatedObject(setting.Path(), processor.Msvm_ProcessorSettingData, &processorSettingData)
}

func (vm *ComputerSystem) MustGetProcessorSettingData() *processor.ProcessorSettingData {
	setting, err := vm.GetProcessorSettingData()
	if err != nil {
		panic(err)
	}
	return setting
}

// GetMemorySettingData returns the memory setting data of the Virtual Machine
func (vm *ComputerSystem) GetMemorySettingData() (*memory.MemorySettingsData, error) {
	setting, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	var memorySettingsData = memory.MemorySettingsData{
		//service: vm.service,
	}
	return &memorySettingsData, vm.GetService().FindFirstRelatedObject(setting.Path(), memory.Msvm_MemorySettingData, &memorySettingsData)
}

func (vm *ComputerSystem) MustGetMemorySettingData() *memory.MemorySettingsData {
	setting, err := vm.GetMemorySettingData()
	if err != nil {
		panic(err)
	}
	return setting
}

const VirtualSystemType_Snapshot = "Microsoft:Hyper-V:Snapshot:Realized"

func (vm *ComputerSystem) GetVirtualSystemSettingData() (*VirtualSystemSettingData, error) {
	var instances []*wmiext.Instance
	var err error
	var virtualSystemSettingData = VirtualSystemSettingData{}
	if instances, err = vm.GetService().FindRelatedInstances(vm.Path(), Msvm_VirtualSystemSettingData); err != nil {
		return nil, err
	}

	var virtualSystemType string
	for _, instance := range instances {
		if virtualSystemType, err = instance.GetAsString("VirtualSystemType"); err != nil {
			return nil, err
		}
		if virtualSystemType == VirtualSystemType_Snapshot {
			continue
		}
		if err = instance.GetAll(&virtualSystemSettingData); err != nil {
			return nil, err
		}
		return &virtualSystemSettingData, nil
	}
	return nil, errors.Wrapf(wmiext.NotFound, "VirtualSystemSettingData not found for computerSystem [%s]", vm.ElementName)
}

func (vm *ComputerSystem) MustGetVirtualSystemSettingData() *VirtualSystemSettingData {
	setting, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		panic(err)
	}
	return setting
}

func (vm *ComputerSystem) Refresh() (err error) {
	instance, err := vm.GetService().GetObject(vm.Path())
	if err != nil {
		return
	}
	if err = instance.Refresh(); err != nil {
		return
	}
	instance, err = vm.GetService().RefetchObject(instance)
	if err != nil {
		return
	}
	return instance.GetAll(vm)
}

func (vm *ComputerSystem) GetState() (ComputerSystemState, error) {
	var err error

	if err = vm.Refresh(); err != nil {
		return Unknown, err
	}
	return ComputerSystemState(vm.EnabledState), nil
}

func (vm *ComputerSystem) GetStatus() (string, error) {
	var err error
	if err = vm.Refresh(); err != nil {
		return "", err
	}
	return vm.Status, nil
}

func (vm *ComputerSystem) GetStatusDescriptions() ([]string, error) {
	err := vm.Refresh()
	if err != nil {
		return nil, err
	}
	return vm.StatusDescriptions, nil
}

func (vm *ComputerSystem) RequestStateChange(requestedState ComputerSystemState) error {
	var err error
	var retValue int32
	var job *wmiext.Instance

	if err = vm.Method("RequestStateChange").
		In("RequestedState", uint16(requestedState)).
		In("TimeoutPeriod", &time.Time{}).
		Execute().
		Out("Job", &job).
		Out("ReturnValue", &retValue).
		End(); err != nil {
		return errors.Wrapf(err, "Failed to request state change to %v", requestedState)
	}

	return utils.WaitResult(retValue, vm.GetService(), job, "Failed to request state change", nil)
}

// ChangeState changes the state of the Virtual Machine
func (vm *ComputerSystem) ChangeState(state ComputerSystemState) (err error) {
	var cState ComputerSystemState
	if cState, err = vm.GetState(); err != nil {
		return
	}
	// If the state is already satisfied, just return
	if cState == state {
		return nil
	}
	// If the state is not satisfied, change the state
	if err = vm.RequestStateChange(state); err != nil {
		return
	}
	return
}

// WaitForState waits for the Virtual Machine to reach the desired state
func (vm *ComputerSystem) WaitForState(state ComputerSystemState, timeoutSeconds int32) (err error) {
	var (
		curState           ComputerSystemState
		vmState            ComputerSystemState
		status             string
		statusDescriptions []string
	)
	start := time.Now()
	// Run the loop, only if the job is actually running
	for {
		if curState, err = vm.GetState(); err != nil {
			return
		} else if curState == state {
			// Break for any valid state
			// TODO: WaitForSomeState
			return nil
		}

		time.Sleep(100 * time.Millisecond)

		// If we have waited enough time, break
		if time.Since(start) > (time.Duration(timeoutSeconds) * time.Second) {
			if vmState, err = vm.GetState(); err != nil {
				vmState = Unknown
			}
			if status, err = vm.GetStatus(); err != nil {
				status = fmt.Sprintf("Unknown (error retreiving the status [%+v])", err)
			}

			if statusDescriptions, err = vm.GetStatusDescriptions(); err != nil {
				statusDescriptions = []string{fmt.Sprintf("Unknown (error retreiving the status descriptions [%+v])", err)}
			}
			err = errors.Wrapf(wmiext.TimedOut, "WaitForState timeout. Current state: [%v], status: [%v], status descriptions: [%v]", vmState, status, statusDescriptions)
			return
		}
	}
}

func (vm *ComputerSystem) RequireState(state ...ComputerSystemState) (ok bool, err error) {
	var curState ComputerSystemState
	if curState, err = vm.GetState(); err != nil {
		return false, err
	}
	for _, s := range state {
		if curState == s {
			return true, nil
		}
	}
	return false, nil
}

// Start starts the Virtual Machine
func (vm *ComputerSystem) Start() (err error) {
	if err = vm.ChangeState(Running); err != nil {
		return
	}
	return vm.WaitForState(Running, StateChangeTimeoutSeconds)
}

// Stop stops the Virtual Machine
func (vm *ComputerSystem) Stop(force bool) (err error) {
	if force {
		if err = vm.ChangeState(Off); err != nil {
			return
		}
	} else {
		if err = vm.ChangeState(Stopping); err != nil {
			// the device is not usable
			if errors.Unwrap(err).(*wmiext.JobError).ErrorCode == 32768 {
				return vm.Stop(force)
			}
			return
		}
	}
	return vm.WaitForState(Off, StateChangeTimeoutSeconds)
}

// ForceStop stops the Virtual Machine immediately
func (vm *ComputerSystem) ForceStop() (err error) {
	return vm.Stop(true)
}

func (vm *ComputerSystem) Reboot(force bool) (err error) {
	var ok bool
	// If the vm is not stopped, stop it
	if ok, err = vm.RequireState(Off); err != nil {
		return
	} else if !ok {
		if err = vm.Stop(force); err != nil {
			return
		}
	}
	// Start the vm
	if err = vm.Start(); err != nil {
		return
	}
	return vm.WaitForState(Running, StateChangeTimeoutSeconds)
}

func (vm *ComputerSystem) ForceReboot() (err error) {
	return vm.Reboot(true)
}

func (vm *ComputerSystem) Pause() (err error) {
	if err = vm.ChangeState(Paused); err != nil {
		return
	}
	return vm.WaitForState(Paused, StateChangeTimeoutSeconds)
}

func (vm *ComputerSystem) Save() (err error) {
	var ok bool
	if ok, err = vm.RequireState(Running); err != nil {
		return
	} else if !ok {
		return errors.New("The virtual machine cannot be saved!")
	}
	if err = vm.ChangeState(Saved); err != nil {
		return
	}
	return vm.WaitForState(Saved, StateChangeTimeoutSeconds)
}

// Resume Virtual Machine
func (vm *ComputerSystem) Resume() error {
	err := vm.ChangeState(Running)
	if err != nil {
		return err
	}
	return vm.WaitForState(Running, StateChangeTimeoutSeconds)
}

// Restore Virtual Machine
func (vm *ComputerSystem) Restore() error {
	err := vm.ChangeState(Running)
	if err != nil {
		return err
	}
	return vm.WaitForState(Running, StateChangeTimeoutSeconds)
}

func (vsms *VirtualSystemManagementService) fetchComputerSystems(wquery string) (
	vms []*ComputerSystem,
	err error,
) {
	err = vsms.Session.FindObjects(wquery, &vms)
	return
}

type VirtualHardDiskType int32

const (
	VirtualHardDiskType_OS_VIRTUALHARDDISK       VirtualHardDiskType = 0
	VirtualHardDiskType_DATADISK_VIRTUALHARDDISK VirtualHardDiskType = 1
)

type HyperVGeneration string

const (
	HyperVGeneration_V1 = "Microsoft:Hyper-V:SubType:1"
	HyperVGeneration_V2 = "Microsoft:Hyper-V:SubType:2"
)

func (vm *ComputerSystem) GetVirtualMachineGeneration() (HyperVGeneration, error) {
	systemSetting, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return "", err
	}
	defer systemSetting.Close()
	value := HyperVGeneration(systemSetting.VirtualSystemSubType)
	return value, nil
}

func (vm *ComputerSystem) GetIDEControllers() ([]*resource.ResourceAllocationSettingData, error) {
	return vm.GetResourceAllocationSettingData(resource.ResourcePool_ResourceType_IDE_Controller)
}

const (
	IDEController  = "IDE Controller"
	SCSIController = "SCSI Controller"
)

func (vm *ComputerSystem) NewSyntheticDiskDrive(
	controllerNumber, controllerLocation int32,
	diskType VirtualHardDiskType,
) (
	synDrive *drive.SyntheticDiskDrive,
	err error,
) {
	// Get the drive allocation pool
	diskDriveResourcePool, err := resource.GetPrimordialResourcePool(vm.GetService(), resource.ResourcePool_ResourceType_Disk_Drive)
	if err != nil {
		return
	}
	defer diskDriveResourcePool.Close()

	resourceAllocationSettingData, err := diskDriveResourcePool.GetDefaultResourceAllocationSettingData()
	if err != nil {
		return
	}

	generation, err := vm.GetVirtualMachineGeneration()
	if err != nil {
		return
	}

	synDrive, err = drive.NewSyntheticDiskDrive(resourceAllocationSettingData.Instance)
	if err != nil {
		return
	}
	//defer synDrive.Close()

	var controllers []*resource.ResourceAllocationSettingData
	var controllerType string

	if generation == HyperVGeneration_V1 && diskType == VirtualHardDiskType_OS_VIRTUALHARDDISK {
		controllers, err = vm.GetIDEControllers()
		controllerType = IDEController
		if err != nil {
			return
		}
	} else {
		controllers, err = vm.GetSCSIControllers()
		controllerType = SCSIController
		if err != nil {
			return
		}
	}

	// 1. Find the correct controller to use vased on the controllerNumber
	if len(controllers) == 0 {
		err = errors.Wrapf(wmiext.NotFound, "VirtualMachine [%s] doesnt have [%s]", vm.ElementName, controllerType)
		return
	}
	if int(controllerNumber) > len(controllers) {
		err = errors.Wrapf(wmiext.NotFound,
			"VirtualMachine [%s] doesnt have [%s] with bus location [%d]", vm.ElementName, controllerType, controllerNumber)
		return
	}

	if controllerNumber == -1 {
		controllerNumber = 0
	}

	var parent string
	var addressOnParent string

	if generation == HyperVGeneration_V1 && diskType == VirtualHardDiskType_OS_VIRTUALHARDDISK {
		ideController := controller.NewIDEControllerSettings(controllers[controllerNumber])
		parent = ideController.Path()

		if controllerLocation == -1 {
			controllerLocation, err = ideController.GetFreeLocation()
			if err != nil {
				err = errors.Wrapf(wmiext.NotFound, "Unable to find free location in IDE Controller")
				return nil, err
			}
			// Find a free location
		}
		addressOnParent = fmt.Sprintf("%d", controllerLocation)
	} else {
		scsiController := controller.NewSCSIControllerSettings(controllers[controllerNumber])
		parent = scsiController.Path()

		if controllerLocation == -1 {
			controllerLocation, err = scsiController.GetFreeLocation()
			if err != nil {
				err = errors.Wrapf(wmiext.NotFound, "Unable to find free location in SCSI Controller")
				return nil, err
			}
			// Find a free location
		}
		addressOnParent = fmt.Sprintf("%d", controllerLocation)
	}

	if err = synDrive.SetParent(parent); err != nil {
		return
	}

	if err = synDrive.SetAddressOnParent(addressOnParent); err != nil {
		return
	}

	return
}

func (vm *ComputerSystem) GetVirtualHardDisks() (col []*disk.VirtualHardDisk, err error) {
	var (
		virtualHardDisk *disk.VirtualHardDisk
	)
	systemSettingData, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return
	}
	storageAllocationSettingDatas, err := systemSettingData.getResourceAllocationSettingData(resource.ResourceAllocationSettingData_ResourceType_Disk_Drive)
	if err != nil {
		return
	}
	for _, storageAllocationSettingData := range storageAllocationSettingDatas {
		if virtualHardDisk, err = disk.NewVirtualHardDisk(storageAllocationSettingData.Instance); err != nil {
			return
		}
		col = append(col, virtualHardDisk)
	}
	return
}

func (vm *ComputerSystem) GetVirtualHardDiskByPath(path string) (targetVirtualHardDisk *disk.VirtualHardDisk, err error) {
	var instance *wmiext.Instance
	virtualHardDisks, err := vm.GetVirtualHardDisks()
	if err != nil {
		return
	}

	for _, virtualHardDisk := range virtualHardDisks {
		virtualHardDiskPath := virtualHardDisk.GetPath()
		if virtualHardDiskPath == path {
			if instance, err = virtualHardDisk.CloneInstance(); err != nil {
				return
			}
			targetVirtualHardDisk, err = disk.NewVirtualHardDisk(instance)
			if err != nil {
				return
			}
			return
		}
	}
	err = errors.Wrapf(wmiext.NotFound, "Vhd with path [%s] not found in Vm [%s]", path, vm.ElementName)
	return
}

func (vm *ComputerSystem) NewVirtualHardDisk(path string) (vhd *disk.VirtualHardDisk, err error) {
	vhdResourcePool, err := resource.GetPrimordialResourcePool(vm.GetService(), resource.ResourcePool_ResourceType_Logical_Disk)
	if err != nil {
		return
	}
	defer vhdResourcePool.Close()
	resourceAllocationSettingData, err := vhdResourcePool.GetDefaultResourceAllocationSettingData()
	if err != nil {
		return
	}

	vhd, err = disk.NewVirtualHardDisk(resourceAllocationSettingData.Instance)
	if err != nil {
		return
	}
	if err = vhd.SetHostResource([]string{path}); err != nil {
		return
	}
	return
}

// NewSyntheticNetworkAdapter creates a new synthetic network adapter for the Virtual Machine, SyntheticNetworkAdapter
func (vm *ComputerSystem) NewSyntheticNetworkAdapter(name string) (adapter *network_adapter.VirtualNetworkAdapter, err error) {
	resourcePool, err := resource.GetPrimordialResourcePool(vm.GetService(), resource.ResourcePool_ResourceType_Ethernet_Adapter)
	if err != nil {
		return
	}
	defer resourcePool.Close()
	rasd, err := resourcePool.GetDefaultResourceAllocationSettingData()
	if err != nil {
		return
	}

	if adapter, err = network_adapter.NewVirtualNetworkAdapterFromInstance(rasd.Instance); err != nil {
		rasd.Close()
		return
	}

	defer func() {
		if err != nil {
			adapter.Close()
			adapter = nil
		}
	}()

	if err = adapter.SetElementName(name); err != nil {
		return
	}

	return
}

// ListComputerSystems returns all the computer systems in the Hyper-V host
func (vsms *VirtualSystemManagementService) ListComputerSystems() ([]*ComputerSystem, error) {
	query := fmt.Sprintf("SELECT * FROM %s", Msvm_ComputerSystem)
	return vsms.fetchComputerSystems(query)
}

// FindComputerSystemsByName returns the computer systems with the given name
func (vsms *VirtualSystemManagementService) FindComputerSystemsByName(name string) ([]*ComputerSystem, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE ElementName='%s'", Msvm_ComputerSystem, name)
	return vsms.fetchComputerSystems(query)
}

// FirstComputerSystemByName returns the computer systems with the given name
func (vsms *VirtualSystemManagementService) FirstComputerSystemByName(name string) (*ComputerSystem, error) {
	computerSystems, err := vsms.FindComputerSystemsByName(name)
	if err != nil {
		return nil, err
	}
	if len(computerSystems) == 0 {
		return nil, wmiext.NotFound
	}
	return computerSystems[0], nil
}

func GetComputerSystem(vna *network_adapter.VirtualNetworkAdapter) (cs *ComputerSystem, err error) {
	instance, err := vna.GetRelated(Msvm_VirtualSystemSettingData)
	if err != nil {
		return
	}
	defer instance.Close()
	virtualSystemSettingData := &VirtualSystemSettingData{}
	if err = instance.GetAll(virtualSystemSettingData); err != nil {
		return
	}
	return virtualSystemSettingData.GetComputerSystem()
}

func (vm *ComputerSystem) NewEthernetPortAllocationSettingData(vna *network_adapter.VirtualNetworkAdapter) (epas *networking.EthernetPortAllocationSettingData, err error) {
	resourcePool, err := resource.GetPrimordialResourcePool(vm.GetService(), resource.ResourcePool_ResourceType_Ethernet_Connection)
	if err != nil {
		return
	}
	defer resourcePool.Close()
	rasd, err := resourcePool.GetDefaultResourceAllocationSettingData()
	if err != nil {
		return
	}
	if err = rasd.SetParent(vna.Path()); err != nil {
		rasd.Close()
		return
	}

	return networking.NewEthernetPortAllocationSettingDataFromInstance(rasd.Instance)
}

func (vm *ComputerSystem) GetVirtualNetworkAdapterByName(name string) (vna *network_adapter.VirtualNetworkAdapter, err error) {
	settings, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return
	}
	defer settings.Close()
	vna, err = settings.GetVirtualNetworkAdapterByName(name)
	return
}
