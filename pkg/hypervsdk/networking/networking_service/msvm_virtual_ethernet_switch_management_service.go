package networking_service

import (
	"fmt"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	utils "github.com/rokukoo/hypervctl/pkg/hypervsdk/utils"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system/host"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"strings"
	"time"
)

const (
	Msvm_VirtualEthernetSwitchManagementService = "Msvm_VirtualEthernetSwitchManagementService"
)

type VirtualEthernetSwitchManagementService struct {
	Con *wmiext.Service
	*wmiext.Instance
}

func (vsms *VirtualEthernetSwitchManagementService) AddResourceSettings(settingsData *networking.VirtualEthernetSwitchSettingData, resourceSettings []*networking.EthernetPortAllocationSettingData) ([]*networking.EthernetPortAllocationSettingData, error) {
	var (
		err                       error
		job                       *wmiext.Instance
		returnValue               int32
		settingsDataRef           string
		resultingResourceSettings []string
	)
	// Get the system settings
	settingsDataRef = settingsData.Path()
	resourceSettingRefs := slice.Map(resourceSettings, func(_ int, rs *networking.EthernetPortAllocationSettingData) string {
		return rs.GetCimText()
	})

	if err = vsms.Method("AddResourceSettings").
		In("AffectedConfiguration", settingsDataRef).
		In("ResourceSettings", resourceSettingRefs).
		Execute().
		Out("Job", &job).
		Out("ResultingResourceSettings", &resultingResourceSettings).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return nil, fmt.Errorf("failed to remove allocation settings: %w", err)
	}
	if err = utils.WaitResult(returnValue, vsms.Con, job, "failed to add allocation settings", nil); err != nil {
		return nil, err
	}
	// Get the allocation settings
	resourceSettings = make([]*networking.EthernetPortAllocationSettingData, len(resultingResourceSettings))
	for i, rs := range resultingResourceSettings {
		target := &networking.EthernetPortAllocationSettingData{}
		if err = vsms.Con.GetObjectAsObject(rs, target); err != nil {
			return nil, err
		}
		resourceSettings[i] = target
	}
	return resourceSettings, nil
}

func (vsms *VirtualEthernetSwitchManagementService) RemoveResourceSettings(resourceSettings []*networking.EthernetPortAllocationSettingData) error {
	var (
		err         error
		job         *wmiext.Instance
		returnValue int32
	)

	// Get the system settings
	resourceSettingPaths := slice.Map(resourceSettings, func(_ int, rs *networking.EthernetPortAllocationSettingData) string {
		return rs.Path()
	})

	for {
		if err = vsms.Method("RemoveResourceSettings").
			In("ResourceSettings", resourceSettingPaths).
			Execute().
			Out("Job", &job).
			Out("ReturnValue", &returnValue).
			End(); err != nil {
			return fmt.Errorf("failed to remove allocation settings: %w", err)
		}
		switch returnValue {
		case 32775:
			// Method failed with 32775, retrying
			time.Sleep(100 * time.Millisecond)
			continue
		case 0:
			return nil
		case 4096:
			if err = utils.WaitResult(returnValue, vsms.Con, job, "failed to remove allocation settings", nil); err != nil {
				return err
			}
			return nil
		}
	}
}

func (vsms *VirtualEthernetSwitchManagementService) DefineSystem(
	settings *networking.VirtualEthernetSwitchSettingData,
	resourceSettings []*networking.EthernetPortAllocationSettingData,
) (*networking.VirtualEthernetSwitch, error) {
	var (
		vswitch = &networking.VirtualEthernetSwitch{}
		err     error

		systemSettings string

		job             *wmiext.Instance
		res             int32
		resultingSystem string

		path            string
		affectedElement *wmiext.Instance
	)
	// Get the system settings
	systemSettings = settings.GetCimText()
	resourceSettingsText := slice.Map(resourceSettings, func(_ int, rs *networking.EthernetPortAllocationSettingData) string {
		return rs.GetCimText()
	})
	// Invoke the DefineSystem method
	err = vsms.Method("DefineSystem").
		In("SystemSettings", systemSettings).
		In("ResourceSettings", resourceSettingsText).
		In("ReferenceConfiguration", nil).
		Execute().
		Out("Job", &job).
		Out("ResultingSystem", &resultingSystem).
		Out("ReturnValue", &res).
		End()
	// Check for errors
	if err != nil {
		return nil, fmt.Errorf("failed to define system: %w", err)
	}
	// Wait for the job to complete
	if err = utils.WaitResult(res, vsms.Con, job, "failed to define system", nil); err != nil {
		return nil, err
	}
	// Get the resulting system
	switch res {
	// Success
	case 0:
		if err = vsms.Con.GetObjectAsObject(resultingSystem, vswitch); err != nil {
			return nil, err
		}
		//if err = vsms.Con.FindFirstRelatedObject(resultingSystem, "Msvm_VirtualEthernetSwitchSettingData", vswitch); err != nil {
		//	return nil, err
		//}
		return vswitch, nil
	// Job in progress
	case 4096:
		// Get the job path
		if path, err = job.Path(); err != nil {
			return nil, err
		}
		// Get the affected element
		if affectedElement, err = vsms.Con.FindFirstRelatedInstanceThrough(path, "Msvm_VirtualEthernetSwitch", "Msvm_AffectedJobElement"); err != nil {
			return nil, err
		}
		// Get the path of the affected element
		if path, err = affectedElement.Path(); err != nil {
			return nil, err
		}
		// Get the virtual switch
		if err = vsms.Con.GetObjectAsObject(path, vswitch); err != nil {
			return nil, err
		}
		return vswitch, nil
	}

	return vswitch, err
}

func (vsms *VirtualEthernetSwitchManagementService) DestroySystem(vswitch *networking.VirtualEthernetSwitch) (err error) {
	var (
		job *wmiext.Instance
		res int32
	)
	for {
		err = vsms.Method("DestroySystem").
			In("AffectedSystem", vswitch.Path()).
			Execute().
			Out("Job", &job).
			Out("ReturnValue", &res).
			End()
		if err != nil {
			return fmt.Errorf("failed to destroy system: %w", err)
		}
		switch res {
		case 32775:
			// Method failed with 32775, retrying
			time.Sleep(100 * time.Millisecond)
			continue
		case 0:
			return nil
		case 4096:
			if err = wmiext.WaitJob(vsms.Con, job); err != nil {
				desc, _ := job.GetAsString("ErrorDescription")
				desc = strings.Replace(desc, "\n", " ", -1)
				return fmt.Errorf("failed to destroy system:%w (%s)", err, desc)
			}
			return nil
		}
	}
}

func (vsms *VirtualEthernetSwitchManagementService) ModifyResourceSettings(elementName string) error {
	return nil
}

func (vsms *VirtualEthernetSwitchManagementService) FirstVirtualSwitchByName(name string) (*networking.VirtualEthernetSwitch, error) {
	var err error
	vswitch := &networking.VirtualEthernetSwitch{}
	wQuery := fmt.Sprintf("SELECT * FROM Msvm_VirtualEthernetSwitch WHERE ElementName = '%s'", name)

	if err = vsms.Con.FindFirstObject(wQuery, vswitch); err != nil {
		return nil, errors.Wrap(err, "failed to find virtual switch")
	}
	return vswitch, nil
}

func LocalVirtualEthernetSwitchManagementService() (*VirtualEthernetSwitchManagementService, error) {
	var (
		con *wmiext.Service
		svc *wmiext.Instance
		err error
	)
	// Get the WMI service
	if con, err = utils.NewLocalHyperVService(); err != nil {
		return nil, err
	}
	// Get the singleton instance
	if svc, err = con.GetSingletonInstance(Msvm_VirtualEthernetSwitchManagementService); err != nil {
		return nil, err
	}
	return &VirtualEthernetSwitchManagementService{con, svc}, nil
}

func MustLocalVirtualEthernetSwitchManagementService() *VirtualEthernetSwitchManagementService {
	vsms, err := LocalVirtualEthernetSwitchManagementService()
	if err != nil {
		panic(err)
	}
	return vsms
}

func (vsms *VirtualEthernetSwitchManagementService) GetVirtualEthernetSwitchSettingData(elementName string) (*networking.VirtualEthernetSwitchSettingData, error) {
	virtualEthernetSwitchSettingData := &networking.VirtualEthernetSwitchSettingData{
		ElementName: elementName,
	}
	instance, err := vsms.Con.CreateInstance(networking.Msvm_VirtualEthernetSwitchSettingData, virtualEthernetSwitchSettingData)
	if err != nil {
		return nil, err
	}
	virtualEthernetSwitchSettingData.Instance = instance
	return virtualEthernetSwitchSettingData, nil
}

func (vsms *VirtualEthernetSwitchManagementService) CreateExternalVirtualSwitch(
	physicalNic,
	externalPortName,
	internalPortName string,
	settings *networking.VirtualEthernetSwitchSettingData,
	enableInternalPort bool,
) (
	vswitch *networking.VirtualEthernetSwitch,
	err error,
) {
	var (
		internalPortAllocSettingData, externalPortAllocSettingData *networking.EthernetPortAllocationSettingData

		resourceSettings []*networking.EthernetPortAllocationSettingData
	)

	// Enable internal port if needed
	if enableInternalPort {
		if internalPortAllocSettingData, err = vsms.DefaultInternalPortAllocationSettingData(internalPortName); err != nil {
			return
		}
		resourceSettings = append(resourceSettings, internalPortAllocSettingData)
	}

	if externalPortAllocSettingData, err = vsms.DefaultExternalPortAllocationSettingData(externalPortName, []string{physicalNic}); err != nil {
		return
	}
	resourceSettings = append(resourceSettings, externalPortAllocSettingData)
	vswitch, err = vsms.DefineSystem(settings, resourceSettings)
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreatePrivateVirtualSwitch(settings *networking.VirtualEthernetSwitchSettingData) (
	vswitch *networking.VirtualEthernetSwitch,
	err error) {
	vswitch, err = vsms.DefineSystem(settings, nil)
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreateInternalVirtualSwitch(internalPortName string, settings *networking.VirtualEthernetSwitchSettingData) (
	vswitch *networking.VirtualEthernetSwitch,
	err error,
) {
	epads, err := vsms.DefaultInternalPortAllocationSettingData(internalPortName)
	if err != nil {
		return
	}
	vswitch, err = vsms.DefineSystem(settings, []*networking.EthernetPortAllocationSettingData{epads})
	return
}

func (vsms *VirtualEthernetSwitchManagementService) GetDefaultEthernetPortAllocationSettingData() (*networking.EthernetPortAllocationSettingData, error) {
	var (
		epasd = &networking.EthernetPortAllocationSettingData{}
		err   error
	)

	instance, err := vsms.Con.CreateInstance(networking.Msvm_EthernetPortAllocationSettingData, epasd)
	if err != nil {
		return nil, err
	}
	epasd.Instance = instance
	return epasd, nil
}

func (vsms *VirtualEthernetSwitchManagementService) DefaultInternalPortAllocationSettingData(switchPortName string) (*networking.EthernetPortAllocationSettingData, error) {
	epasd, err := vsms.GetDefaultEthernetPortAllocationSettingData()

	hostCm, err := host.GetHostComputerSystem()
	if err != nil {
		return nil, err
	}
	defer hostCm.Close()

	epasd.HostResource = []string{hostCm.Path()}
	epasd.ElementName = switchPortName

	return epasd, epasd.PutAll(epasd)
}

func (vsms *VirtualEthernetSwitchManagementService) DefaultExternalPortAllocationSettingData(
	switchPortName string,
	physicalNicNames []string,
) (
	*networking.EthernetPortAllocationSettingData,
	error,
) {
	epasd, err := vsms.GetDefaultEthernetPortAllocationSettingData()
	if err != nil {
		return nil, err
	}

	if len(physicalNicNames) == 0 {
		return nil, errors.Wrapf(err, "Physical Nic Name is missing")
	}

	epasd.ElementName = switchPortName

	hresource := []string{}
	for _, nicName := range physicalNicNames {
		// Get the External Ethernet Port
		if extPort, err := networking.GetExternalEthernetPort(vsms.Con, nicName); err != nil {
			// If the External Ethernet Port is not found, try to get the WiFi Port
			if errors.Is(err, wmiext.NotFound) {
				wifiPort, err := networking.GetWiFiPort(vsms.Con, nicName)
				if err != nil {
					return nil, err
				}
				epasd.Address = wifiPort.PermanentAddress
				hresource = append(hresource, wifiPort.Path())
			} else {
				return nil, err
			}
		} else {
			epasd.Address = extPort.PermanentAddress
			hresource = append(hresource, extPort.Path())
		}
	}

	epasd.HostResource = hresource

	return epasd, epasd.PutAll(epasd)
}

func (vsms *VirtualEthernetSwitchManagementService) ClearInternalPortAllocationSettingData(vsw *networking.VirtualEthernetSwitch) error {
	settings, err := vsw.GetInternalPortAllocSettings()
	if err != nil {
		if errors.Is(err, wmiext.NotFound) {
			return nil
		}
		return errors.Wrap(err, "failed to get internal port allocation settings")
	}
	return vsms.RemoveResourceSettings([]*networking.EthernetPortAllocationSettingData{settings})
}

func (vsms *VirtualEthernetSwitchManagementService) ClearExternalPortAllocationSettingData(vsw *networking.VirtualEthernetSwitch) error {
	settings, err := vsw.GetExternalPortAllocSettings()
	if err != nil {
		if errors.Is(err, wmiext.NotFound) {
			return nil
		}
		return errors.Wrap(err, "failed to get internal port allocation settings")
	}
	return vsms.RemoveResourceSettings([]*networking.EthernetPortAllocationSettingData{settings})
}
