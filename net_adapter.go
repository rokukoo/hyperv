package hypervctl

import (
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/errors"
	errors2 "github.com/microsoft/wmi/pkg/errors"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	"github.com/microsoft/wmi/pkg/virtualization/network/switchport"
	na "github.com/microsoft/wmi/pkg/virtualization/network/virtualnetworkadapter"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	"github.com/rokukoo/hypervctl/wmictl"
)

// VirtualNetworkAdapter represents a virtual network adapter
// VirtualMachine: The hyper-v virtual machine
// VirtualNetworkAdapter: The hyper-v virtual network adapter, any operation on the virtual network adapter is done through this
type VirtualNetworkAdapter struct {
	VirtualMachine *HyperVVirtualMachine
	*na.VirtualNetworkAdapter
}

// Name returns the name of the virtual network adapter
func (vna *VirtualNetworkAdapter) Name() (string, error) {
	name, err := vna.VirtualNetworkAdapter.GetPropertyElementName()
	if err != nil {
		return "", err
	}
	return name, nil
}

// Connect connects the virtual network adapter to a virtual switch
func (vna *VirtualNetworkAdapter) Connect(vsw *VirtualSwitch) (bool, error) {
	var (
		err            error
		vmms           *service.VirtualSystemManagementService
		virtualMachine *virtualsystem.VirtualMachine
		name           string
	)
	// Get the virtual system management service
	if vmms, err = wmictl.NewLocalVirtualSystemManagementService(); err != nil {
		return false, err
	}
	// Get the hyper-v physical virtual machine
	if virtualMachine, err = vna.VirtualMachine.VM(); err != nil {
		return false, err
	}
	// Get the name of the virtual network adapter
	if name, err = vna.Name(); err != nil {
		return false, err
	}
	// Connect the virtual network adapter to the virtual switch
	if err = vmms.ConnectAdapterToVirtualSwitch(virtualMachine, name, vsw.VirtualSwitch); err != nil {
		return false, err
	}
	return true, nil
}

// ModifyConfiguration modifies the network adapter configuration
// ipV4Address: The IPv4 address
// subnetMask: The subnet mask
// defaultGateway: The default gateway
func (vna *VirtualNetworkAdapter) ModifyConfiguration(ipV4Address []string, subnetMask []string, defaultGateway []string, dnsServer []string) (bool, error) {
	var (
		err error
	)
	// Get the hyper-v virtual network adapter
	virtualNetworkAdapter := vna.VirtualNetworkAdapter
	// Get the network adapter guest configuration,
	// which is provided by the hyper-v and supports modifying the network adapter configuration without connecting to the virtual machine
	configuration, err := virtualNetworkAdapter.GetGuestNetworkAdapterConfiguration()
	if err != nil {
		return false, err
	}
	if err = configuration.SetPropertyDHCPEnabled(false); err != nil {
		return false, err
	}
	if err = configuration.SetProperty("ProtocolIFType", 4096); err != nil {
		return false, err
	}
	if err = configuration.SetPropertyIPAddresses(ipV4Address); err != nil {
		return false, err
	}
	if err = configuration.SetPropertySubnets(subnetMask); err != nil {
		return false, err
	}
	if err = configuration.SetPropertyDefaultGateways(defaultGateway); err != nil {
		return false, err
	}
	if err = configuration.SetPropertyDNSServers(dnsServer); err != nil {
		return false, err
	}
	// Apply the network adapter configuration
	virtualMachine, err := vna.VirtualMachine.VM()
	if err != nil {
		return false, err
	}
	vmms, err := wmictl.NewLocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	vmInstancePath := virtualMachine.InstancePath()
	embeddedConfigurationInstance, err := configuration.EmbeddedXMLInstance()
	if err != nil {
		return false, err
	}
	// Get the method to set the network adapter configuration
	method, err := vmms.GetWmiMethod("SetGuestNetworkAdapterConfiguration")
	if err != nil {
		return false, err
	}
	// Execute the method to set the network adapter configuration
	inparams := wmi.WmiMethodParamCollection{}
	inparams = append(inparams, wmi.NewWmiMethodParam("ComputerSystem", vmInstancePath))
	inparams = append(inparams, wmi.NewWmiMethodParam("NetworkConfiguration", []string{embeddedConfigurationInstance}))
	outparams := wmi.WmiMethodParamCollection{wmi.NewWmiMethodParam("Job", nil)}
	outparams = append(outparams, wmi.NewWmiMethodParam("ReturnValue", nil))
	// Execute the method to set the network adapter configuration
	result, err := method.Execute(inparams, outparams)
	if err != nil {
		return false, err
	}

	if !(result.ReturnValue == 4096 || result.ReturnValue == 0) {
		err = errors.Wrapf(errors2.Failed, "Method failed with [%d]", result.ReturnValue)
		return false, nil
	}

	if result.ReturnValue == 0 {
		return true, nil
	}

	val, ok := result.OutMethodParams["Job"]
	if !ok || val.Value == nil {
		err = errors.Wrapf(errors2.NotFound, "Job")
		return false, err
	}
	job, err := instance.GetWmiJob(vmms.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
	if err != nil {
		return false, err
	}
	defer job.Close()
	err = job.WaitForJobCompletion(result.ReturnValue, -1)
	return true, err
}

// SetBandwidth sets the bandwidth of the virtual network adapter
// limitMbps: The maximum bandwidth in Mbps
// reserveMbps: The minimum bandwidth in Mbps
func (vna *VirtualNetworkAdapter) SetBandwidth(limitMbps float64, reserveMbps float64) error {
	if limitMbps < 0 {
		return errors.New("limitMbps must be greater than or equal to 0")
	}
	vsms, err := wmictl.GetVirtualSystemManagementService(vna.VirtualMachine.GetWmiHost())
	if err != nil {
		return err
	}
	// Get the virtual network adapter
	syntheticAdapter := vna.VirtualNetworkAdapter
	ethernetPortAllocationSettingData, err := syntheticAdapter.GetEthernetPortAllocationSettingData()
	// If the virtual network adapter does not contain an ethernet port allocation setting data, create a new one
	if ethernetPortAllocationSettingData == nil {
		virtualMachine, err := vna.VirtualMachine.VM()
		if err != nil {
			return err
		}
		ethernetPortAllocationSettingData, err = vsms.AddVirtualEthernetConnection(virtualMachine, syntheticAdapter)
		if err != nil {
			return err
		}
		if ethernetPortAllocationSettingData == nil {
			return errors.New("Failed to add virtual ethernet connection")
		}
	}
	// Get the virtual network adapter bandwidth setting data
	related, err := ethernetPortAllocationSettingData.GetRelated("Msvm_EthernetSwitchPortBandwidthSettingData")
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	// If the virtual network adapter bandwidth setting data does not exist, create a new one
	if related == nil {
		// Create a new virtual network adapter bandwidth setting data
		bandwidthSettingData, err := DefaultEthernetSwitchPortBandwidthSettingData(vsms)
		if err != nil {
			return err
		}
		if err = bandwidthSettingData.SetPropertyLimit(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		//if err = bandwidthSettingData.SetPropertyReservation(uint64(reserveMbps * 100000)); err != nil {
		//	return err
		//}
		if err = bandwidthSettingData.SetPropertyBurstLimit(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		if err = bandwidthSettingData.SetPropertyBurstSize(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		err = vsms.AddEthernetFeatureEx1(ethernetPortAllocationSettingData.Msvm_EthernetPortAllocationSettingData, bandwidthSettingData.WmiInstance, -1)
		return err
	} else {
		// Modify the existing virtual network adapter bandwidth setting data
		bandwidthSettingData, err := switchport.NewEthernetSwitchPortBandwidthSettingData(related)
		if err != nil {
			return err
		}
		defer bandwidthSettingData.Close()
		if err = bandwidthSettingData.SetPropertyLimit(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		//if err = bandwidthSettingData.SetPropertyReservation(uint64(reserveMbps * 1000000)); err != nil {
		//	return err
		//}
		if err = bandwidthSettingData.SetPropertyBurstLimit(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		if err = bandwidthSettingData.SetPropertyBurstSize(uint64(limitMbps * 1000000)); err != nil {
			return err
		}
		_, err = vsms.ModifyEthernetFeature(wmi.WmiInstanceCollection{bandwidthSettingData.WmiInstance}, -1)
		return err
	}
}

// DefaultEthernetSwitchPortBandwidthSettingData returns the default EthernetSwitchPortBandwidthSettingData
func DefaultEthernetSwitchPortBandwidthSettingData(vsms *service.VirtualSystemManagementService) (*switchport.EthernetSwitchPortBandwidthSettingData, error) {
	hc, err := vsms.GetHostComputerSystem()
	if err != nil {
		return nil, err
	}
	tmp, err := hc.GetDefaultPortSettingData("Ethernet Switch Port Bandwidth Settings", "Msvm_EthernetSwitchPortBandwidthSettingData")
	if err != nil {
		return nil, err
	}
	spbs, err := switchport.NewEthernetSwitchPortBandwidthSettingData(tmp)
	if err != nil {
		return nil, err
	}
	return spbs, nil
}

// virtualNetworkAdapter creates a new virtual network adapter entity from a hyper-v virtual network adapter
func virtualNetworkAdapter(vna *na.VirtualNetworkAdapter) (*VirtualNetworkAdapter, error) {
	return &VirtualNetworkAdapter{
		VirtualNetworkAdapter: vna,
	}, nil
}
