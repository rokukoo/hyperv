package wmiext

import (
	"encoding/xml"
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/errors"
	diskService "github.com/microsoft/wmi/pkg/virtualization/core/storage/service"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	v2 "github.com/microsoft/wmi/server2019/root/virtualization/v2"
	"log"
	"time"
)

func newLocalImageManagementService(whost *host.WmiHost) (mgmt *diskService.ImageManagementService, err error) {
	creds := whost.GetCredential()
	query := query.NewWmiQuery("Msvm_ImageManagementService")
	// TODO: Regenerate wrappers that would take WmiHost directly
	imswmi, err := v2.NewMsvm_ImageManagementServiceEx6(whost.HostName, string(constant.Virtualization), creds.UserName, creds.Password, creds.Domain, query)
	if err != nil {
		return
	}
	mgmt = &diskService.ImageManagementService{imswmi}
	return
}

func NewLocalImageManagementService() (*diskService.ImageManagementService, error) {
	whost := host.NewWmiLocalHost()
	mgmt, err := newLocalImageManagementService(whost)
	if err != nil {
		return nil, err
	}
	return mgmt, nil
}

func GetVirtualHardDiskState(path string) (virtualHardDiskState *VirtualHardDiskState, err error) {
	ims, err := NewLocalImageManagementService()
	if err != nil {
		return
	}

	method, err := ims.GetWmiMethod("GetVirtualHardDiskState")
	if err != nil {
		return
	}

	inparams := wmi.WmiMethodParamCollection{}
	inparams = append(inparams, wmi.NewWmiMethodParam("Path", path))
	outparams := wmi.WmiMethodParamCollection{wmi.NewWmiMethodParam("State", nil)}
	outparams = append(outparams, wmi.NewWmiMethodParam("Job", nil))

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

		// Try to get the Out Params
		state := result.OutMethodParams["State"]
		if state != nil {
			ins := &virtualHardDiskStateInstance{}
			err = xml.Unmarshal([]byte(state.Value.(string)), ins)
			if err != nil {
				return
			}
			virtualHardDiskState, err = ins.virtualHardDiskState()
			if err != nil {
				return
			}
		}

		if result.ReturnValue == 0 {
			return
		}

		val, ok := result.OutMethodParams["Job"]
		if !ok || val.Value == nil {
			err = errors.Wrapf(errors.NotFound, "Job")
			return
		}
		job, err1 := instance.GetWmiJob(ims.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
		if err1 != nil {
			err = err1
			return
		}
		defer job.Close()
		err = job.WaitForJobCompletion(result.ReturnValue, -1)
		return
	}
	return
}
