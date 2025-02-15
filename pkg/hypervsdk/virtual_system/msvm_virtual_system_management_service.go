package virtual_system

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/memory"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/network_adapter"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/processor"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/resource"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/disk"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/storage/drive"
	utils "github.com/rokukoo/hypervctl/pkg/hypervsdk/utils"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"strings"
	"time"
)

const (
	Msvm_VirtualSystemManagementService = "Msvm_VirtualSystemManagementService"
)

type VirtualSystemManagementService struct {
	Session *wmiext.Service
	*wmiext.Instance
}

func LocalVirtualSystemManagementService() (*VirtualSystemManagementService, error) {
	var (
		session *wmiext.Service
		svc     *wmiext.Instance
		err     error
	)
	// Get the WMI service
	if session, err = utils.NewLocalHyperVService(); err != nil {
		return nil, err
	}
	// Get the singleton instance
	if svc, err = session.GetSingletonInstance(Msvm_VirtualSystemManagementService); err != nil {
		return nil, err
	}
	return &VirtualSystemManagementService{session, svc}, nil
}

func MustLocalVirtualSystemManagementService() *VirtualSystemManagementService {
	vsms, err := LocalVirtualSystemManagementService()
	if err != nil {
		panic(err)
	}
	return vsms
}

func (vsms *VirtualSystemManagementService) DefineSystem(
	systemSettingsData *VirtualSystemSettingData,
	processorSetting *processor.ProcessorSettingData,
	memorySetting *memory.MemorySettingsData,
) (*ComputerSystem, error) {
	var (
		system = ComputerSystem{}

		err                   error
		systemSettingsDataObj string
		memorySettingObj      string
		processorSettingObj   string

		job             *wmiext.Instance
		returnValue     int32
		resultingSystem string
	)

	if systemSettingsDataObj, err = vsms.CreateSystemSettings(systemSettingsData); err != nil {
		return nil, err
	}

	if memorySettingObj, err = memory.CreateMemorySettings(memorySetting); err != nil {
		return nil, err
	}

	if processorSettingObj, err = processor.CreateProcessorSettings(processorSetting); err != nil {
		return nil, err
	}

	if err = vsms.Method("DefineSystem").
		In("SystemSettings", systemSettingsDataObj).
		In("ResourceSettings", []string{memorySettingObj, processorSettingObj}).
		In("ReferenceConfiguration", nil).
		Execute().
		Out("Job", &job).
		Out("ResultingSystem", &resultingSystem).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return nil, err
	}

	if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to define system", nil); err != nil {
		return nil, err
	}

	return &system, vsms.Session.GetObjectAsObject(resultingSystem, &system)
}

func (vsms *VirtualSystemManagementService) ModifyProcessorSettings(
	processorSettingData *processor.ProcessorSettingData,
) (err error) {
	_, err = vsms.ModifyResourceSettings([]string{processorSettingData.GetCimText()})
	return
}

func (vsms *VirtualSystemManagementService) ModifyMemorySettings(
	memorySettingData *memory.MemorySettingsData,
) (err error) {
	_, err = vsms.ModifyResourceSettings([]string{memorySettingData.GetCimText()})
	return
}

func (vsms *VirtualSystemManagementService) DestroySystem(
	computerSystem *ComputerSystem,
) error {
	var (
		err error

		job         *wmiext.Instance
		returnValue int32
	)

	for {
		if err = vsms.Method("DestroySystem").
			In("AffectedSystem", computerSystem.Path()).
			Execute().
			Out("Job", &job).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return err
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to destroy system", nil); err != nil {
			return err
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return nil
	}
}

// AddResourceSettings - 将资源添加到虚拟机配置。
//
// 当应用于“状态”虚拟机配置时，作为一种副作用，资源会添加到活动虚拟机。
//
// Microsoft Docs: https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/addresourcesettings-msvm-virtualsystemmanagementservice
func (vsms *VirtualSystemManagementService) AddResourceSettings(
	affectedConfiguration *VirtualSystemSettingData,
	resourceSettings []string,
) (
	[]*wmiext.Instance,
	error,
) {
	var (
		err error

		job                       *wmiext.Instance
		returnValue               int32
		resultingResourceSettings []string

		resultInstances []*wmiext.Instance
	)

	for {
		if err = vsms.Method("AddResourceSettings").
			In("AffectedConfiguration", affectedConfiguration.Path()).
			In("ResourceSettings", resourceSettings).
			Execute().
			Out("Job", &job).
			Out("ResultingResourceSettings", &resultingResourceSettings).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return resultInstances, err
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to add allocation settings", nil); err != nil {
			return resultInstances, err
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, resourceSetting := range resultingResourceSettings {
			var instance *wmiext.Instance
			if instance, err = vsms.Session.GetObject(resourceSetting); err != nil {
				return resultInstances, err
			}
			resultInstances = append(resultInstances, instance)
		}

		return resultInstances, nil
	}
}

func (vsms *VirtualSystemManagementService) ModifyResourceSettings(
	resourceSettings []string,
) (
	resultInstances []*wmiext.Instance,
	err error,
) {
	var (
		job                       *wmiext.Instance
		returnValue               int32
		resultingResourceSettings []string
	)

	for {
		if err = vsms.Method("ModifyResourceSettings").
			In("ResourceSettings", resourceSettings).
			Execute().
			Out("Job", &job).
			Out("ResultingResourceSettings", &resultingResourceSettings).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return resultInstances, err
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to modify allocation settings", nil); err != nil {
			return resultInstances, err
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, resourceSetting := range resultingResourceSettings {
			var instance *wmiext.Instance
			if instance, err = vsms.Session.GetObject(resourceSetting); err != nil {
				return resultInstances, err
			}
			resultInstances = append(resultInstances, instance)
		}

		return resultInstances, nil
	}
}

func (vsms *VirtualSystemManagementService) AddSCSIController(
	vm *ComputerSystem,
) (
	err error,
) {
	var (
		result         []*wmiext.Instance
		scsiController *resource.ResourceAllocationSettingData
		systemSetting  *VirtualSystemSettingData
	)

	if scsiController, err = vm.NewSCSIController(); err != nil {
		return err
	}

	if systemSetting, err = vm.GetVirtualSystemSettingData(); err != nil {
		return
	}

	// apply the settings
	if result, err = vsms.AddResourceSettings(systemSetting, []string{scsiController.GetCimText()}); err != nil {
		return
	}

	if len(result) == 0 {
		err = errors.Wrapf(wmiext.NotFound, "AddVirtualSystemResource")
		return
	}

	return

}

func (vsms *VirtualSystemManagementService) RemoveResourceSettings(
	resourceSettings []string,
) (
	err error,
) {
	var (
		job         *wmiext.Instance
		returnValue int32
	)

	for {
		if err = vsms.Method("RemoveResourceSettings").
			In("ResourceSettings", resourceSettings).
			Execute().
			Out("Job", &job).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to remove allocation settings", nil); err != nil {
			return
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return nil
	}
}

func (vsms *VirtualSystemManagementService) RemoveSyntheticDiskDrive(diskDrive *drive.SyntheticDiskDrive) error {
	return vsms.RemoveResourceSettings([]string{diskDrive.Path()})
}

func (vsms *VirtualSystemManagementService) AddSyntheticDiskDrive(
	vm *ComputerSystem,
	controllerNumber, controllerLocation int32,
	diskType VirtualHardDiskType,
) (
	vhdDrive *drive.SyntheticDiskDrive,
	err error,
) {
	systemSettingData, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return
	}
	defer systemSettingData.Close()

	syntheticDiskDrive, err := vm.NewSyntheticDiskDrive(controllerNumber, controllerLocation, diskType)
	if err != nil {
		return
	}
	defer syntheticDiskDrive.Close()
	resultInstances, err := vsms.AddResourceSettings(systemSettingData, []string{syntheticDiskDrive.GetCimText()})
	if err != nil {
		return
	}
	if len(resultInstances) == 0 {
		err = errors.Wrapf(wmiext.NotFound, "AddVirtualSystemResource")
		return
	}
	driveInstance, err := resultInstances[0].CloneInstance()
	if err != nil {
		return
	}

	vhdDrive, err = drive.NewSyntheticDiskDrive(driveInstance)
	return
}

// AttachVirtualHardDisk -
// * Create a Synthetic Disk Drive
// *    Add a drive to available first controller at available location
// * Connects the Disk to the Drive
// Returns Disk and Drive
func (vsms *VirtualSystemManagementService) AttachVirtualHardDisk(
	vm *ComputerSystem,
	path string,
	diskType VirtualHardDiskType,
) (
	vhd *disk.VirtualHardDisk,
	vhdDrive *drive.SyntheticDiskDrive,
	err error,
) {

	// Add a drive
	if vhdDrive, err = vsms.AddSyntheticDiskDrive(vm, -1, -1, diskType); err != nil {
		return
	}

	defer func() {
		if err != nil {
			err1 := vsms.RemoveSyntheticDiskDrive(vhdDrive)
			if err1 != nil {
				//log.Printf("RemoveSyntheticDiskDrive [%+v]\n", err1)
			}
			vhdDrive.Close()
			vhdDrive = nil
		}
	}()

	// Add a disk
	virtualHardDisk, err := vm.NewVirtualHardDisk(path)
	if err != nil {
		return
	}
	defer virtualHardDisk.Close()

	// ConnectByName disk to drive
	if err = virtualHardDisk.SetParent(vhdDrive.Path()); err != nil {
		return
	}

	if !strings.Contains(virtualHardDisk.Path(), "Definition") {
		var resultInstances []*wmiext.Instance
		if resultInstances, err = vsms.ModifyResourceSettings([]string{virtualHardDisk.GetCimText()}); err != nil {
			return
		}
		if len(resultInstances) == 0 {
			err = errors.Wrapf(wmiext.NotFound, "AddVirtualSystemResource")
			return
		}
		if virtualHardDisk, err = disk.NewVirtualHardDisk(resultInstances[0]); err != nil {
			return
		}
	}

	systemSettingData, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return
	}
	defer systemSettingData.Close()

	// apply the settings
	resultInstances, err := vsms.AddResourceSettings(systemSettingData, []string{virtualHardDisk.GetCimText()})
	if err != nil {
		return
	}

	if len(resultInstances) == 0 {
		// Sometimes this could hapen - Find out why
		vhd, err = vm.GetVirtualHardDiskByPath(path)
		return
	}

	vhdInstance, err := resultInstances[0].CloneInstance()
	if err != nil {
		return
	}

	vhd, err = disk.NewVirtualHardDisk(vhdInstance)
	if err != nil {
		vhdInstance.Close()
	}
	return
}

func (vsms *VirtualSystemManagementService) DetachVirtualHardDisk(virtualHardDisk *disk.VirtualHardDisk) (err error) {
	diskDrive, err := virtualHardDisk.GetDrive()
	if err != nil {
		return
	}
	defer diskDrive.Close()

	// Remove Disk
	if err = vsms.RemoveResourceSettings([]string{virtualHardDisk.Path()}); err != nil {
		return
	}
	// Remove Drive
	if err = vsms.RemoveResourceSettings([]string{diskDrive.Path()}); err != nil {
		return
	}
	return
}

// AddFeatureSettings - Add feature settings to a virtual system
//
// Microsoft Documentation: https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/addfeaturesettings-msvm-virtualsystemmanagementservice
func (vsms *VirtualSystemManagementService) AddFeatureSettings(
	affectedConfiguration string,
	featureSettings []string,
) (
	resultInstances []*wmiext.Instance,
	err error,
) {
	var (
		resultingFeatureSettings []string

		job         *wmiext.Instance
		returnValue int32
	)

	for {
		if err = vsms.Method("AddFeatureSettings").
			In("AffectedConfiguration", affectedConfiguration).
			In("FeatureSettings", featureSettings).
			Execute().
			Out("Job", &job).
			Out("ResultingFeatureSettings", &resultingFeatureSettings).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to modify allocation settings", nil); err != nil {
			return
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, resourceSetting := range resultingFeatureSettings {
			var instance *wmiext.Instance
			if instance, err = vsms.Session.GetObject(resourceSetting); err != nil {
				return resultInstances, err
			}
			resultInstances = append(resultInstances, instance)
		}

		return resultInstances, nil

	}
}

// ModifyFeatureSettings - 修改虚拟机以太网连接的当前功能设置。
//
// Microsoft Documentation: https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/modifyfeaturesettings-msvm-virtualsystemmanagementservice
func (vsms *VirtualSystemManagementService) ModifyFeatureSettings(
	featureSettings []string,
) (
	resultInstances []*wmiext.Instance,
	err error,
) {
	var (
		resultingFeatureSettings []string

		job         *wmiext.Instance
		returnValue int32
	)

	for {
		if err = vsms.Method("ModifyFeatureSettings").
			In("FeatureSettings", featureSettings).
			Execute().
			Out("Job", &job).
			Out("ResultingFeatureSettings", &resultingFeatureSettings).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to modify allocation settings", nil); err != nil {
			return
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		for _, resourceSetting := range resultingFeatureSettings {
			var inst *wmiext.Instance
			if inst, err = vsms.Session.GetObject(resourceSetting); err != nil {
				return resultInstances, err
			}
			resultInstances = append(resultInstances, inst)
		}

		return resultInstances, nil

	}
}

// RemoveFeatureSettings - 删除虚拟机以太网连接的当前功能设置
//
// Microsoft Documentation: https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/removefeaturesettings-msvm-virtualsystemmanagementservice
func (vsms *VirtualSystemManagementService) RemoveFeatureSettings(
	featureSettings []string,
) (
	err error,
) {
	var (
		job         *wmiext.Instance
		returnValue int32
	)

	for {
		if err = vsms.Method("RemoveFeatureSettings").
			In("FeatureSettings", featureSettings).
			Execute().
			Out("Job", &job).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return
		}

		if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to modify allocation settings", nil); err != nil {
			return
		}

		if returnValue == 32775 {
			//log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", "DestroySystem", returnValue)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		return nil
	}
}

func (vsms *VirtualSystemManagementService) DisConnectAdapterToVirtualSwitch(vnaName string) (err error) {
	adapter, err := network_adapter.FirstVirtualNetworkAdapterByName(vsms.Session, vnaName)
	if err != nil {
		return
	}

	// If the adapter is already connected, then Msvm_EthernetPortAllocationSettingData would exists, if not it would be null
	pasd, err := adapter.GetEthernetPortAllocationSettingData()
	if err != nil {
		if errors.Is(err, wmiext.NotFound) {
			return nil
		}
		return
	}
	defer pasd.Close()
	err = vsms.RemoveResourceSettings([]string{pasd.Path()})
	return
}

func (vsms *VirtualSystemManagementService) ConnectAdapterToVirtualSwitch(computerSystem *ComputerSystem, vnaName string, vsw *networking.VirtualEthernetSwitch) (err error) {
	adapter, err := computerSystem.GetVirtualNetworkAdapterByName(vnaName)
	if err != nil {
		return
	}
	defer adapter.Close()

	// If the adapter is already connected, then Msvm_EthernetPortAllocationSettingData would exists, if not it would be null
	pasd, err := adapter.GetEthernetPortAllocationSettingData()
	if err != nil {
		if !errors.Is(err, wmiext.NotFound) {
			return
		}
		pasd, err = vsms.AddVirtualEthernetConnection(computerSystem, adapter)
		if err != nil {
			return
		}
	} else {
		// Already existing case
		if err = pasd.SetEnabledState(2); err != nil {
			return
		}
	}
	defer pasd.Close()
	if err = pasd.SetHostResource([]string{vsw.Path()}); err != nil {
		return
	}
	_, err = vsms.ModifyResourceSettings([]string{pasd.GetCimText()})
	return
}

// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/setguestnetworkadapterconfiguration-msvm-virtualsystemmanagementservice
func (vsms *VirtualSystemManagementService) SetGuestNetworkAdapterConfiguration(
	computerSystem *ComputerSystem,
	networkConfiguration *network_adapter.GuestNetworkAdapterConfiguration,
) (
	err error,
) {
	var (
		job         *wmiext.Instance
		returnValue int32
	)

	if err = vsms.Method("SetGuestNetworkAdapterConfiguration").
		In("computerSystem", computerSystem.Path()).
		In("NetworkConfiguration", []string{networkConfiguration.GetCimText()}).
		Execute().
		Out("Job", &job).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return
	}

	if err = utils.WaitResult(returnValue, vsms.Session, job, "Failed to modify allocation settings", nil); err != nil {
		return
	}

	return
}
