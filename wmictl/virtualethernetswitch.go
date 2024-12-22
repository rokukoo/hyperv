package wmictl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/errors"
	"github.com/microsoft/wmi/pkg/hardware/network/netadapter"
	netservice "github.com/microsoft/wmi/pkg/hardware/network/service"
	vmmshost "github.com/microsoft/wmi/pkg/virtualization/core/host"
	"github.com/microsoft/wmi/pkg/virtualization/network/ethernetport"
	"github.com/microsoft/wmi/pkg/virtualization/network/switchport"
	"github.com/microsoft/wmi/pkg/virtualization/network/virtualswitch"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	v2 "github.com/microsoft/wmi/server2019/root/virtualization/v2"
)

// "github.com/microsoft/wmi/pkg/virtualization/network/virtualethernetswitch"

func GetWifiPort(whost *host.WmiHost, ethernetName string) (wifiport *v2.Msvm_WiFiPort, err error) {
	creds := whost.GetCredential()
	query := query.NewWmiQuery("Msvm_WiFiPort", "ElementName", ethernetName)
	wifiport, err = v2.NewMsvm_WiFiPortEx6(whost.HostName, string(constant.Virtualization), creds.UserName, creds.Password, creds.Domain, query)
	if err != nil {
		return
	}
	return
}

func GetExternalPortAllocationSettingData(whost *host.WmiHost, switchPortName string, physicalNicNames []string) (epas *switchport.EthernetPortAllocationSettingData, err error) {
	// Get the default Ethernet Port Allocation Setting Data
	if epas, err = vmmshost.GetDefaultEthernetPortAllocationSettingData(whost); err != nil {
		return nil, err
	}

	if len(physicalNicNames) == 0 {
		err = errors.Wrapf(errors.InvalidInput, "Physical Nic Name is missing")
		return
	}
	err = epas.SetPropertyElementName(switchPortName)
	if err != nil {
		return
	}

	hresource := []string{}
	for _, nicName := range physicalNicNames {
		// Get the External Ethernet Port
		if extPort, err := ethernetport.GetExternalEthernetPort(whost, nicName); err != nil {
			// If the External Ethernet Port is not found, try to get the WiFi Port
			if errors.IsNotFound(err) {
				wifiPort, err := GetWifiPort(whost, nicName)
				if err != nil {
					return nil, err
				}
				//address, err := wifiPort.GetPropertyPermanentAddress()
				//if err != nil {
				//	return nil, err
				//}
				//err = epas.SetPropertyAddress(address)
				//if err != nil {
				//	return nil, err
				//}
				hresource = append(hresource, wifiPort.InstancePath())
			} else {
				return nil, err
			}
		} else {
			address, err := extPort.GetPropertyPermanentAddress()
			if err != nil {
				return nil, err
			}
			err = epas.SetPropertyAddress(address)
			if err != nil {
				return nil, err
			}
			hresource = append(hresource, extPort.InstancePath())
		}
	}

	err = epas.SetPropertyHostResource(hresource)
	if err != nil {
		return
	}
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreateExternalVirtualSwitch(physicalNicName, externalPortName, internalPortName string,
	settings *virtualswitch.VirtualEthernetSwitchSettingData, internalport bool) (
	vswitch *virtualswitch.VirtualSwitch,
	err error) {
	var (
		adapter              *netadapter.NetworkAdapter
		interfaceDescription string

		internalpasd *switchport.EthernetPortAllocationSettingData
		externalpasd *switchport.EthernetPortAllocationSettingData

		internalpasdXml string
		externalpasdXml string

		resourceSettings []string
	)

	if internalport {
		if internalpasd, err = vmmshost.GetInternalPortAllocationSettingData(vsms.GetWmiHost(), internalPortName); err != nil {
			return
		}
		if internalpasdXml, err = internalpasd.EmbeddedXMLInstance(); err != nil {
			return
		}
		resourceSettings = append(resourceSettings, internalpasdXml)
	}

	if adapter, err = netservice.GetNetworkAdapterByName(vsms.GetWmiHost(), physicalNicName); err != nil {
		return
	}
	if interfaceDescription, err = adapter.GetPropertyInterfaceDescription(); err != nil {
		return
	}
	if externalpasd, err = GetExternalPortAllocationSettingData(vsms.GetWmiHost(), externalPortName, []string{interfaceDescription}); err != nil {
		return
	}

	if externalpasdXml, err = externalpasd.EmbeddedXMLInstance(); err != nil {
		return
	}
	resourceSettings = append(resourceSettings, externalpasdXml)
	vswitch, err = vsms.CreateVirtualSwitch(settings, resourceSettings)
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreatePrivateVirtualSwitch(settings *virtualswitch.VirtualEthernetSwitchSettingData) (
	vswitch *virtualswitch.VirtualSwitch,
	err error) {
	vswitch, err = vsms.CreateVirtualSwitch(settings, []string{})
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreateInternalVirtualSwitch(internalPortName string, settings *virtualswitch.VirtualEthernetSwitchSettingData) (
	vswitch *virtualswitch.VirtualSwitch,
	err error) {
	internalpasd, err := vmmshost.GetInternalPortAllocationSettingData(vsms.GetWmiHost(), internalPortName)
	if err != nil {
		return
	}
	defer internalpasd.Close()
	internalpasdXml, err := internalpasd.EmbeddedXMLInstance()
	if err != nil {
		return
	}

	vswitch, err = vsms.CreateVirtualSwitch(settings, []string{internalpasdXml})
	return
}

func (vsms *VirtualEthernetSwitchManagementService) CreateVirtualSwitch(settings *virtualswitch.VirtualEthernetSwitchSettingData, resourceSettings []string) (
	vswitch *virtualswitch.VirtualSwitch,
	err error) {

	method, err := vsms.GetWmiMethod("DefineSystem")
	if err != nil {
		return
	}
	defer method.Close()

	embeddedInstance, err := settings.EmbeddedXMLInstance()
	if err != nil {
		return
	}

	inparams := wmi.WmiMethodParamCollection{}
	inparams = append(inparams, wmi.NewWmiMethodParam("SystemSettings", embeddedInstance))
	inparams = append(inparams, wmi.NewWmiMethodParam("ResourceSettings", resourceSettings))
	inparams = append(inparams, wmi.NewWmiMethodParam("ReferenceConfiguration", nil))
	outparams := wmi.WmiMethodParamCollection{wmi.NewWmiMethodParam("Job", nil)}
	outparams = append(outparams, wmi.NewWmiMethodParam("ResultingSystem", nil))

	result, err := method.Execute(inparams, outparams)
	if err != nil {
		return
	}

	if !(result.ReturnValue == 4096 || result.ReturnValue == 0) {
		err = errors.Wrapf(errors.Failed, "Method failed with [%d]", result.ReturnValue)
		return
	}
	val, ok := result.OutMethodParams["ResultingSystem"]
	if ok && val.Value != nil {
		vswitchinstance, err := instance.GetWmiInstanceFromPath(vsms.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
		if err == nil {
			vswitch, err = virtualswitch.NewVirtualSwitch(vswitchinstance)
		}
	}

	if result.ReturnValue == 0 {
		return
	}

	val, ok = result.OutMethodParams["Job"]
	if !ok || val.Value == nil {
		err = errors.Wrapf(errors.NotFound, "Job")
		return
	}
	job, err := instance.GetWmiJob(vsms.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
	if err != nil {
		return
	}
	defer job.Close()
	err = job.WaitForJobCompletion(result.ReturnValue, -1)

	if vswitch != nil {
		return
	}

	affectedElement, err := job.GetFirstRelatedEx("Msvm_AffectedJobElement", "", "", "")
	if err != nil {
		// For now, ignore the error
		err = nil
		return
	}
	vswitch, err = virtualswitch.NewVirtualSwitch(affectedElement)
	return
}
