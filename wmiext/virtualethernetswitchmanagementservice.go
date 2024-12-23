package wmiext

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/errors"
	"github.com/microsoft/wmi/pkg/virtualization/network/virtualswitch"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	v2 "github.com/microsoft/wmi/server2019/root/virtualization/v2"
	"log"
	"time"
)

type VirtualEthernetSwitchManagementService struct {
	*v2.Msvm_VirtualEthernetSwitchManagementService
}

// GetVirtualSwitches would get all virtual machines
func (vsms *VirtualEthernetSwitchManagementService) GetVirtualSwitches() (*virtualswitch.VirtualSwitchCollection, error) {
	wHost := host.NewWmiLocalHost()
	query := query.NewWmiQuery("Msvm_VirtualEthernetSwitch")
	instances, err := instance.GetWmiInstancesFromHost(wHost, VirtualizationV2, query)
	if err != nil {
		return nil, err
	}
	vmc := virtualswitch.VirtualSwitchCollection{}
	for _, ins := range instances {
		vm, err := virtualswitch.NewVirtualSwitch(ins)
		if err != nil {
			return nil, err
		}
		vmc = append(vmc, vm)
	}
	return &vmc, nil
}

func (vsms *VirtualEthernetSwitchManagementService) FindVirtualSwitchByName(vswitchName string) (*virtualswitch.VirtualSwitch, error) {
	vswitchs, err := vsms.GetVirtualSwitches()
	if err != nil {
		return nil, err
	}

	for _, entity := range *vswitchs {
		entityName, err := entity.GetPropertyElementName()
		if err != nil {
			return nil, err
		}
		if entityName != vswitchName {
			continue
		}

		clonedEntity, err := entity.Clone()
		if err != nil {
			return nil, err
		}
		return virtualswitch.NewVirtualSwitch(clonedEntity)
	}

	return nil, errors.Wrapf(errors.NotFound, "Unable to find a virtual system with name [%s]", vswitchName)
}

func (vsms *VirtualEthernetSwitchManagementService) DeleteVirtualSwitch(vswitch *virtualswitch.VirtualSwitch) (err error) {
	method, err := vsms.GetWmiMethod("DestroySystem")
	if err != nil {
		return
	}
	defer method.Close()

	inparams := wmi.WmiMethodParamCollection{}
	inparams = append(inparams, wmi.NewWmiMethodParam("AffectedSystem", vswitch.InstancePath()))
	outparams := wmi.WmiMethodParamCollection{wmi.NewWmiMethodParam("Job", nil)}

	for {
		result, err1 := method.Execute(inparams, outparams)
		if err1 != nil {
			err = err1
			return
		}

		returnVal := result.ReturnValue
		if returnVal != 0 && returnVal != 4096 {
			// Virtual System is in Invalid State, try to retry
			if returnVal == 32775 {
				log.Printf("[WMI] Method [%s] failed with [%d]. Retrying ...", method.Name, returnVal)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			err = errors.Wrapf(errors.Failed, "Method failed with [%d]", result.ReturnValue)
			return
		}

		if result.ReturnValue == 0 {
			return
		}

		val, ok := result.OutMethodParams["Job"]
		if !ok || val.Value == nil {
			err = errors.Wrapf(errors.NotFound, "Job")
			return
		}
		job, err1 := instance.GetWmiJob(vsms.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
		if err1 != nil {
			err = err1
			return
		}
		defer job.Close()
		return job.WaitForJobCompletion(result.ReturnValue, -1)
	}
	return
}

func newLocalVirtualEthernetSwitchManagementService(whost *host.WmiHost) (mgmt *VirtualEthernetSwitchManagementService, err error) {
	creds := whost.GetCredential()
	query := query.NewWmiQuery("Msvm_VirtualEthernetSwitchManagementService")
	// TODO: Regenerate wrappers that would take WmiHost directly
	vswitchwmi, err := v2.NewMsvm_VirtualEthernetSwitchManagementServiceEx6(whost.HostName, string(constant.Virtualization), creds.UserName, creds.Password, creds.Domain, query)
	if err != nil {
		return
	}

	mgmt = &VirtualEthernetSwitchManagementService{vswitchwmi}
	return
}

func NewLocalVirtualEthernetSwitchManagementService() (*VirtualEthernetSwitchManagementService, error) {
	whost := host.NewWmiLocalHost()
	mgmt, err := newLocalVirtualEthernetSwitchManagementService(whost)
	if err != nil {
		return nil, err
	}
	return mgmt, nil
}
