package hypervctl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/network_adapter"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking/switch_extension"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system/host"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualNetworkAdapter struct {
	Name        string
	Description string

	IsEnableVlan bool
	VlanId       int

	MacAddress string

	Ipv4Address []string
	Gateway     []string
	SubnetMask  []string
	DnsServer   []string

	IsEnableBandwidth bool
	MinBandwidth      float64
	MaxBandwidth      float64
	*network_adapter.VirtualNetworkAdapter
}

func NewVirtualNetworkAdapter(networkAdapter *network_adapter.VirtualNetworkAdapter) (*VirtualNetworkAdapter, error) {
	vna := &VirtualNetworkAdapter{}
	return vna, vna.update(networkAdapter)
}

func (vm *VirtualMachine) FirstVirtualNetworkAdapterByName(name string) (*VirtualNetworkAdapter, error) {
	virtualSystemSettingData, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	instances, err := virtualSystemSettingData.GetAllRelated(networking.Msvm_SyntheticEthernetPortSettingData)
	if err != nil {
		return nil, err
	}
	for _, instance := range instances {
		var networkAdapter *network_adapter.VirtualNetworkAdapter
		if networkAdapter, err = network_adapter.NewVirtualNetworkAdapterFromInstance(instance); err != nil {
			return nil, err
		}
		var vna *VirtualNetworkAdapter
		vna, err = NewVirtualNetworkAdapter(networkAdapter)
		if err != nil {
			return nil, err
		}
		if vna.ElementName != name {
			continue
		}
		return vna, nil
	}
	return nil, wmiext.NotFound
}

func (vm *VirtualMachine) MustFirstVirtualNetworkAdapterByName(name string) *VirtualNetworkAdapter {
	vna, err := vm.FirstVirtualNetworkAdapterByName(name)
	if err != nil {
		panic(err)
	}
	return vna
}

func (vm *VirtualMachine) FindVirtualNetworkAdapterByName(name string) ([]*VirtualNetworkAdapter, error) {
	virtualSystemSettingData, err := vm.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	instances, err := virtualSystemSettingData.GetAllRelated(networking.Msvm_SyntheticEthernetPortSettingData)
	if err != nil {
		return nil, err
	}
	var virtualNetworkAdapters []*VirtualNetworkAdapter
	for _, instance := range instances {
		var networkAdapter *network_adapter.VirtualNetworkAdapter
		if networkAdapter, err = network_adapter.NewVirtualNetworkAdapterFromInstance(instance); err != nil {
			return nil, err
		}
		var vna *VirtualNetworkAdapter
		vna, err = NewVirtualNetworkAdapter(networkAdapter)
		if err != nil {
			return nil, err
		}
		if vna.ElementName != name {
			continue
		}
		virtualNetworkAdapters = append(virtualNetworkAdapters, vna)
	}
	return virtualNetworkAdapters, nil
}

func (vna *VirtualNetworkAdapter) GetVirtualMachine() (*VirtualMachine, error) {
	cs, err := virtual_system.GetComputerSystem(vna.VirtualNetworkAdapter)
	if err != nil {
		return nil, err
	}
	return NewVirtualMachine(cs)
}

func (vna *VirtualNetworkAdapter) update(syntheticNetworkAdapter *network_adapter.VirtualNetworkAdapter) error {
	vna.VirtualNetworkAdapter = syntheticNetworkAdapter
	vna.Name = syntheticNetworkAdapter.ElementName
	return nil
}

func (vm *VirtualMachine) AddVirtualNetworkAdapter(vna *VirtualNetworkAdapter) (err error) {
	var syntheticNetworkAdapter *network_adapter.VirtualNetworkAdapter
	if syntheticNetworkAdapter, err = vm.NewSyntheticNetworkAdapter(vna.Name); err != nil {
		return
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	resourceSettings, err := vmms.AddResourceSettings(vm.MustGetVirtualSystemSettingData(), []string{syntheticNetworkAdapter.GetCimText()})
	if err != nil {
		return
	}
	if len(resourceSettings) == 0 {
		return errors.New("Failed to add resource settings")
	}
	if syntheticNetworkAdapter, err = network_adapter.NewVirtualNetworkAdapterFromInstance(resourceSettings[0]); err != nil {
		return
	}

	vna.VirtualNetworkAdapter = syntheticNetworkAdapter

	if vna.IsEnableBandwidth {
		if err = vna.SetBandwidthOut(vna.MaxBandwidth, vna.MinBandwidth); err != nil {
			return
		}
	}

	if err = vna.update(syntheticNetworkAdapter); err != nil {
		return
	}
	return nil
}

func (vna *VirtualNetworkAdapter) Detach() error {
	if vna.VirtualNetworkAdapter == nil {
		return errors.New("vna not attached")
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	if err = vmms.RemoveResourceSettings([]string{vna.VirtualNetworkAdapter.Path()}); err != nil {
		return err
	}
	vna.VirtualNetworkAdapter = nil
	return nil
}

func (vm *VirtualMachine) RemoveVirtualNetworkAdapter(name string) (err error) {
	networkAdapters, err := vm.FindVirtualNetworkAdapterByName(name)
	if err != nil {
		return
	}
	for _, networkAdapter := range networkAdapters {
		if err = networkAdapter.Detach(); err != nil {
			return
		}
	}
	return
}

// SetBandwidthOut sets the bandwidth of the virtual network adapter
// limitBandwidthMbps: The maximum bandwidth in Mbps, -1 means unlimited
// reserveBandwidthMbps: The minimum bandwidth in Mbps -1 means unlimited
func (vna *VirtualNetworkAdapter) SetBandwidthOut(limitBandwidthMbps, reserveBandwidthMbps float64) (err error) {
	if limitBandwidthMbps < 0 {
		limitBandwidthMbps = 0
	}
	if reserveBandwidthMbps < 0 {
		reserveBandwidthMbps = 0
	}
	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	// Get the virtual network adapter
	syntheticAdapter := vna.VirtualNetworkAdapter
	ethernetPortAllocationSettingData, err := syntheticAdapter.GetEthernetPortAllocationSettingData()
	// If the virtual network adapter does not contain an ethernet port allocation setting data, create a new one
	virtualMachine, err := vna.GetVirtualMachine()
	if err != nil {
		return err
	}
	if ethernetPortAllocationSettingData == nil {
		if ethernetPortAllocationSettingData, err = vsms.AddVirtualEthernetConnection(virtualMachine.ComputerSystem, syntheticAdapter); err != nil {
			return
		}
		if ethernetPortAllocationSettingData == nil {
			return errors.New("Failed to add virtual ethernet connection")
		}
	}

	ethernetSwitchPortBandwidthSettingData, err := ethernetPortAllocationSettingData.GetRelated(switch_extension.Msvm_EthernetSwitchPortBandwidthSettingData)
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return err
	}

	var bandwidthSettingData *switch_extension.EthernetSwitchPortBandwidthSettingData

	// If the virtual network adapter bandwidth setting data does not exist, create a new one
	if ethernetSwitchPortBandwidthSettingData == nil {
		// Build a new virtual network adapter bandwidth setting data
		if bandwidthSettingData, err = host.DefaultEthernetSwitchPortBandwidthSettingData(); err != nil {
			return err
		}
		defer bandwidthSettingData.Close()

		if err = bandwidthSettingData.SetLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
			return
		}
		if err = bandwidthSettingData.SetReservation(uint64(reserveBandwidthMbps * 100000)); err != nil {
			return err
		}
		//if err = bandwidthSettingData.SetBurstLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		//if err = bandwidthSettingData.SetBurstSize(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		_, err = vsms.AddFeatureSettings(ethernetPortAllocationSettingData.Path(), []string{bandwidthSettingData.GetCimText()})
		return
	} else {
		// Modify the existing virtual network adapter bandwidth setting data
		bandwidthSettingData, err = switch_extension.NewEthernetSwitchPortBandwidthSettingData(ethernetSwitchPortBandwidthSettingData)
		if err != nil {
			return err
		}
		defer bandwidthSettingData.Close()
		if err = bandwidthSettingData.SetLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
			return
		}
		if err = bandwidthSettingData.SetReservation(uint64(reserveBandwidthMbps * 100000)); err != nil {
			return err
		}
		//if err = bandwidthSettingData.SetBurstLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		//if err = bandwidthSettingData.SetBurstSize(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		_, err = vsms.ModifyResourceSettings([]string{bandwidthSettingData.GetCimText()})
		return
	}
}

// Connect connects the virtual network adapter to a virtual switch
func (vna *VirtualNetworkAdapter) Connect(vswName string) (bool, error) {
	var (
		err            error
		vsms           *virtual_system.VirtualSystemManagementService
		virtualMachine *VirtualMachine
		vsw            *networking.VirtualEthernetSwitch
	)
	// Get the virtual system management service
	if vsms, err = virtual_system.LocalVirtualSystemManagementService(); err != nil {
		return false, err
	}

	if virtualMachine, err = vna.GetVirtualMachine(); err != nil {
		return false, err
	}

	vnaName := vna.Name
	// Connect the virtual network adapter to the virtual switch
	if vsw, err = networking.FirstVirtualEthernetSwitchByName(vsms.Session, vswName); err != nil {
		return false, err
	}
	if err = vsms.ConnectAdapterToVirtualSwitch(virtualMachine.ComputerSystem, vnaName, vsw); err != nil {
		return false, err
	}
	return true, nil
}

// DisConnect disconnects the virtual network adapter from a virtual switch
func (vna *VirtualNetworkAdapter) DisConnect() (err error) {
	if err = virtual_system.MustLocalVirtualSystemManagementService().DisConnectAdapterToVirtualSwitch(vna.Name); err != nil {
		return
	}
	return
}

//
//// ModifyConfiguration modifies the network adapter configuration
//// ipV4Address: The IPv4 address
//// subnetMask: The subnet mask
//// defaultGateway: The default gateway
//func (vna *VirtualNetworkAdapter) ModifyConfiguration(ipV4Address []string, subnetMask []string, defaultGateway []string, dnsServer []string) (bool, error) {
//	var (
//		err error
//	)
//	// Get the hyper-v virtual network adapter
//	virtualNetworkAdapter := vna.VirtualNetworkAdapter
//	// Get the network adapter guest configuration,
//	// which is provided by the hyper-v and supports modifying the network adapter configuration without connecting to the virtual machine
//	configuration, err := virtualNetworkAdapter.GetGuestNetworkAdapterConfiguration()
//	if err != nil {
//		return false, err
//	}
//	if err = configuration.SetPropertyDHCPEnabled(false); err != nil {
//		return false, err
//	}
//	if err = configuration.SetProperty("ProtocolIFType", 4096); err != nil {
//		return false, err
//	}
//	if err = configuration.SetPropertyIPAddresses(ipV4Address); err != nil {
//		return false, err
//	}
//	if err = configuration.SetPropertySubnets(subnetMask); err != nil {
//		return false, err
//	}
//	if err = configuration.SetPropertyDefaultGateways(defaultGateway); err != nil {
//		return false, err
//	}
//	if err = configuration.SetPropertyDNSServers(dnsServer); err != nil {
//		return false, err
//	}
//	// Apply the network adapter configuration
//	virtualMachine, err := vna.VirtualMachine.VM()
//	if err != nil {
//		return false, err
//	}
//	vmms, err := wmiext.NewLocalVirtualSystemManagementService()
//	if err != nil {
//		return false, err
//	}
//	vmInstancePath := virtualMachine.InstancePath()
//	embeddedConfigurationInstance, err := configuration.EmbeddedXMLInstance()
//	if err != nil {
//		return false, err
//	}
//	// Get the method to set the network adapter configuration
//	method, err := vmms.GetWmiMethod("SetGuestNetworkAdapterConfiguration")
//	if err != nil {
//		return false, err
//	}
//	// Execute the method to set the network adapter configuration
//	inparams := wmi.WmiMethodParamCollection{}
//	inparams = append(inparams, wmi.NewWmiMethodParam("ComputerSystem", vmInstancePath))
//	inparams = append(inparams, wmi.NewWmiMethodParam("NetworkConfiguration", []string{embeddedConfigurationInstance}))
//	outparams := wmi.WmiMethodParamCollection{wmi.NewWmiMethodParam("Job", nil)}
//	outparams = append(outparams, wmi.NewWmiMethodParam("ReturnValue", nil))
//	// Execute the method to set the network adapter configuration
//	result, err := method.Execute(inparams, outparams)
//	if err != nil {
//		return false, err
//	}
//
//	if !(result.ReturnValue == 4096 || result.ReturnValue == 0) {
//		err = errors.Wrapf(errors2.Failed, "Method failed with [%d]", result.ReturnValue)
//		return false, nil
//	}
//
//	if result.ReturnValue == 0 {
//		return true, nil
//	}
//
//	val, ok := result.OutMethodParams["Job"]
//	if !ok || val.Value == nil {
//		err = errors.Wrapf(errors2.NotFound, "Job")
//		return false, err
//	}
//	job, err := instance.GetWmiJob(vmms.GetWmiHost(), string(constant.Virtualization), val.Value.(string))
//	if err != nil {
//		return false, err
//	}
//	defer job.Close()
//	err = job.WaitForJobCompletion(result.ReturnValue, -1)
//	return true, err
//}

// FindVirtualNetworkAdapterByName returns the virtual network adapter
func FindVirtualNetworkAdapterByName(name string) (virtualNetworkAdapters []*VirtualNetworkAdapter, err error) {
	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	wquery := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s'", name)
	return virtualNetworkAdapters, vsms.Session.FindObjects(wquery, virtualNetworkAdapters)
}

// FirstVirtualNetworkAdapterByName returns the first virtual network adapter
func FirstVirtualNetworkAdapterByName(name string) (virtualNetworkAdapter *VirtualNetworkAdapter, err error) {
	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	wquery := fmt.Sprintf("SELECT * FROM Msvm_SyntheticEthernetPortSettingData WHERE ElementName = '%s'", name)
	instance, err := vsms.Session.FindFirstInstance(wquery)
	if err != nil {
		return nil, err
	}
	networkAdapter, err := network_adapter.NewVirtualNetworkAdapterFromInstance(instance)
	if err != nil {
		return nil, err
	}
	return NewVirtualNetworkAdapter(networkAdapter)
}

func MustFirstVirtualNetworkAdapterByName(name string) *VirtualNetworkAdapter {
	vna, err := FirstVirtualNetworkAdapterByName(name)
	if err != nil {
		panic(err)
	}
	return vna
}
