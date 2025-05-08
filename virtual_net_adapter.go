package hypervctl

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/network_adapter"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system"
	"github.com/rokukoo/hypervctl/pkg/hypervsdk/virtual_system/host"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

type VirtualNetworkAdapter struct {
	Name string

	StaticMacAddress bool
	MacAddress       string

	IPAddress      []string
	DefaultGateway []string
	SubnetMask     []string
	DNSServers     []string

	IsEnableBandwidth bool
	MinBandwidth      float64
	MaxBandwidth      float64

	IsEnableVlan bool
	VlanId       int

	virtualNetworkAdapter *network_adapter.VirtualNetworkAdapter
}

func (vm *VirtualMachine) GetVirtualNetworkAdapters() ([]*VirtualNetworkAdapter, error) {
	var virtualNetworkAdapters []*VirtualNetworkAdapter
	virtualSystemSettingData, err := vm.computerSystem.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	syntheticVirtualNetworkAdapters, err := virtualSystemSettingData.GetSyntheticVirtualNetworkAdapters()
	if err != nil {
		return nil, err
	}
	for _, syntheticVirtualNetworkAdapter := range syntheticVirtualNetworkAdapters {
		var vna *VirtualNetworkAdapter
		if vna, err = NewVirtualNetworkAdapter(syntheticVirtualNetworkAdapter); err != nil {
			return nil, err
		}
		virtualNetworkAdapters = append(virtualNetworkAdapters, vna)
	}
	return virtualNetworkAdapters, nil
}

func (vna *VirtualNetworkAdapter) GetVirtualMachine() (*VirtualMachine, error) {
	cs, err := virtual_system.GetComputerSystem(vna.virtualNetworkAdapter)
	if err != nil {
		return nil, err
	}
	return NewVirtualMachine(cs)
}

func (vna *VirtualNetworkAdapter) update(virtualNetworkAdapter *network_adapter.VirtualNetworkAdapter) (err error) {
	vna.virtualNetworkAdapter = virtualNetworkAdapter
	vna.Name = virtualNetworkAdapter.ElementName
	syntheticNetworkAdapter, err := network_adapter.NewSyntheticNetworkAdapterFromInstance(virtualNetworkAdapter.Instance)
	if err != nil {
		return
	}
	vna.StaticMacAddress = syntheticNetworkAdapter.StaticMacAddress
	vna.MacAddress = syntheticNetworkAdapter.Address

	ethernetPortAllocationSettingData, err := syntheticNetworkAdapter.GetEthernetPortAllocationSettingData()
	if err != nil {
		return
	}
	// Bandwidth setting data
	ethernetSwitchPortBandwidthSettingData, err := ethernetPortAllocationSettingData.GetEthernetSwitchPortBandwidthSettingData()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return
	}
	if ethernetSwitchPortBandwidthSettingData != nil {
		vna.IsEnableBandwidth = true
		vna.MaxBandwidth = float64(ethernetSwitchPortBandwidthSettingData.Limit) / 1000000
		vna.MinBandwidth = float64(ethernetSwitchPortBandwidthSettingData.Reservation) / 1000000
	}

	// Vlan setting data
	ethernetSwitchPortVlanSettingData, err := ethernetPortAllocationSettingData.GetEthernetSwitchPortVlanSettingData()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return
	}
	if ethernetSwitchPortVlanSettingData != nil {
		vna.IsEnableVlan = true
		vna.VlanId = int(ethernetSwitchPortVlanSettingData.AccessVlanId)
	}

	configuration, err := vna.virtualNetworkAdapter.GetGuestNetworkAdapterConfiguration()
	if err != nil {
		return
	}
	vna.IPAddress = configuration.IPAddresses
	vna.DefaultGateway = configuration.DefaultGateways
	vna.SubnetMask = configuration.Subnets
	vna.DNSServers = configuration.DNSServers
	return nil
}

func NewVirtualNetworkAdapter(networkAdapter *network_adapter.VirtualNetworkAdapter) (*VirtualNetworkAdapter, error) {
	vna := &VirtualNetworkAdapter{}
	return vna, vna.update(networkAdapter)
}

func (vm *VirtualMachine) FirstVirtualNetworkAdapterByName(name string) (*VirtualNetworkAdapter, error) {
	virtualSystemSettingData, err := vm.computerSystem.GetVirtualSystemSettingData()
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
		if vna.virtualNetworkAdapter.ElementName != name {
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
	virtualSystemSettingData, err := vm.computerSystem.GetVirtualSystemSettingData()
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
		if vna.virtualNetworkAdapter.ElementName != name {
			continue
		}
		virtualNetworkAdapters = append(virtualNetworkAdapters, vna)
	}
	return virtualNetworkAdapters, nil
}

func (vm *VirtualMachine) AddVirtualNetworkAdapter(vna *VirtualNetworkAdapter) (err error) {
	var syntheticNetworkAdapter *network_adapter.VirtualNetworkAdapter
	if syntheticNetworkAdapter, err = vm.computerSystem.NewSyntheticNetworkAdapter(vna.Name); err != nil {
		return
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	resourceSettings, err := vmms.AddResourceSettings(vm.computerSystem.MustGetVirtualSystemSettingData(), []string{syntheticNetworkAdapter.GetCimText()})
	if err != nil {
		return
	}
	if len(resourceSettings) == 0 {
		return errors.New("Failed to add resource settings")
	}
	if syntheticNetworkAdapter, err = network_adapter.NewVirtualNetworkAdapterFromInstance(resourceSettings[0]); err != nil {
		return
	}

	vna.virtualNetworkAdapter = syntheticNetworkAdapter

	if vna.IsEnableBandwidth {
		if err = vna.SetBandwidth(vna.MaxBandwidth, vna.MinBandwidth); err != nil {
			return
		}
	}

	if err = vna.update(syntheticNetworkAdapter); err != nil {
		return
	}
	return nil
}

func (vna *VirtualNetworkAdapter) Detach() error {
	if vna.virtualNetworkAdapter == nil {
		return errors.New("vna not attached")
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	if err = vmms.RemoveResourceSettings([]string{vna.virtualNetworkAdapter.Path()}); err != nil {
		return err
	}
	vna.virtualNetworkAdapter = nil
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

func (vna *VirtualNetworkAdapter) DisableBandwidthLimit() (err error) {
	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	// Get the virtual network adapter
	syntheticAdapter := vna.virtualNetworkAdapter
	ethernetPortAllocationSettingData, err := syntheticAdapter.GetEthernetPortAllocationSettingData()
	if err != nil {
		return err
	}

	ethernetSwitchPortBandwidthSettingData, err := ethernetPortAllocationSettingData.GetEthernetSwitchPortBandwidthSettingData()
	if err != nil {
		if errors.Is(err, wmiext.NotFound) {
			return nil
		}
		return
	}

	err = vsms.RemoveFeatureSettings([]string{ethernetSwitchPortBandwidthSettingData.Path()})

	return
}

// SetBandwidth sets the bandwidth of the virtual network adapter
// limitBandwidthMbps: The maximum bandwidth in Mbps, -1 means unlimited
// reserveBandwidthMbps: The minimum bandwidth in Mbps -1 means unlimited
func (vna *VirtualNetworkAdapter) SetBandwidth(limitBandwidthMbps, reserveBandwidthMbps float64) (err error) {
	if limitBandwidthMbps < 0 {
		limitBandwidthMbps = 0
	}
	if reserveBandwidthMbps < 0 {
		reserveBandwidthMbps = 0
	}
	if limitBandwidthMbps <= reserveBandwidthMbps {
		return wmiext.NotSupported
	}
	if limitBandwidthMbps == 0 && reserveBandwidthMbps == 0 {
		return wmiext.NotSupported
	}
	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return err
	}
	// Get the virtual network adapter
	syntheticAdapter := vna.virtualNetworkAdapter
	ethernetPortAllocationSettingData, _ := syntheticAdapter.GetEthernetPortAllocationSettingData()
	// If the virtual network adapter does not contain an ethernet port allocation setting data, create a new one
	virtualMachine, err := vna.GetVirtualMachine()
	if err != nil {
		return err
	}
	if ethernetPortAllocationSettingData == nil {
		if ethernetPortAllocationSettingData, err = vsms.AddVirtualEthernetConnection(virtualMachine.computerSystem, syntheticAdapter); err != nil {
			return
		}
		if ethernetPortAllocationSettingData == nil {
			return errors.New("Failed to add virtual ethernet connection")
		}
	}

	ethernetSwitchPortBandwidthSettingData, err := ethernetPortAllocationSettingData.GetEthernetSwitchPortBandwidthSettingData()
	if err != nil && !errors.Is(err, wmiext.NotFound) {
		return err
	}
	//defer ethernetSwitchPortBandwidthSettingData.Close()

	modifyBandwidthSettingData := func() error {
		// Set the maximum bandwidth of the virtual network adapter
		if err = ethernetSwitchPortBandwidthSettingData.SetLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
			return err
		}
		// Set the minimum bandwidth of the virtual network adapter
		if err = ethernetSwitchPortBandwidthSettingData.SetReservation(uint64(reserveBandwidthMbps * 1000000)); err != nil {
			return err
		}
		//if err = ethernetSwitchPortBandwidthSettingData.SetBurstLimit(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		//if err = ethernetSwitchPortBandwidthSettingData.SetBurstSize(uint64(limitBandwidthMbps * 1000000)); err != nil {
		//	return err
		//}
		return nil
	}

	// If the virtual network adapter bandwidth setting data does not exist, create a new one
	if ethernetSwitchPortBandwidthSettingData == nil {
		// Build a new virtual network adapter bandwidth setting data
		if ethernetSwitchPortBandwidthSettingData, err = host.DefaultEthernetSwitchPortBandwidthSettingData(); err != nil {
			return err
		}
		if err = modifyBandwidthSettingData(); err != nil {
			return
		}
		_, err = vsms.AddFeatureSettings(ethernetPortAllocationSettingData.Path(), []string{ethernetSwitchPortBandwidthSettingData.GetCimText()})
	} else {
		// ModifySpec the existing virtual network adapter bandwidth setting data
		if err = modifyBandwidthSettingData(); err != nil {
			return
		}
		_, err = vsms.ModifyFeatureSettings([]string{ethernetSwitchPortBandwidthSettingData.GetCimText()})
	}

	// Manually update the bandwidth setting data
	vna.MaxBandwidth = limitBandwidthMbps
	vna.MinBandwidth = reserveBandwidthMbps

	return
}

var (
	ErrorNotConnected = errors.New("vna not connected to virtual switch")
)

func (vna *VirtualNetworkAdapter) GetVirtualSwitch() (*VirtualSwitch, error) {
	// Get the virtual network adapter
	syntheticAdapter := vna.virtualNetworkAdapter
	ethernetPortAllocationSettingData, err := syntheticAdapter.GetEthernetPortAllocationSettingData()
	if err != nil {
		return nil, err
	}
	if ethernetPortAllocationSettingData == nil {
		return nil, ErrorNotConnected
	}
	hostResource := ethernetPortAllocationSettingData.HostResource
	if len(hostResource) == 0 {
		return nil, ErrorNotConnected
	}
	vswPath := hostResource[0]
	virtualEthernetSwitch := &networking.VirtualEthernetSwitch{}
	if err = syntheticAdapter.GetService().GetObjectAsObject(vswPath, virtualEthernetSwitch); err != nil {
		return nil, err
	}
	return NewVirtualSwitch(virtualEthernetSwitch)
}

// ConnectByName connects the virtual network adapter to a virtual switch
func (vna *VirtualNetworkAdapter) ConnectByName(vswName string) (bool, error) {
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
	// ConnectByName the virtual network adapter to the virtual switch
	if vsw, err = networking.FirstVirtualEthernetSwitchByName(vsms.Session, vswName); err != nil {
		return false, err
	}
	if err = vsms.ConnectAdapterToVirtualSwitch(virtualMachine.computerSystem, vnaName, vsw); err != nil {
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

// ModifyConfiguration modifies the network adapter configuration
// ipV4Address: The IPv4 address
// subnetMask: The subnet mask
// defaultGateway: The default gateway
func (vna *VirtualNetworkAdapter) ModifyConfiguration(
	ipV4Address, subnetMask, defaultGateway, dnsServer []string,
) (err error) {
	var guestNetworkAdapterConfiguration *network_adapter.GuestNetworkAdapterConfiguration
	// Get the network adapter guest guestNetworkAdapterConfiguration,
	// which is provided by the hyper-v and supports modifying the network adapter guestNetworkAdapterConfiguration without connecting to the virtual machine
	if guestNetworkAdapterConfiguration, err = vna.virtualNetworkAdapter.GetGuestNetworkAdapterConfiguration(); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetDHCPEnabled(false); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetProtocolIFType(network_adapter.IPv4); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetIPAddresses(ipV4Address); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetSubnets(subnetMask); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetDefaultGateways(defaultGateway); err != nil {
		return
	}
	if err = guestNetworkAdapterConfiguration.SetDNSServers(dnsServer); err != nil {
		return
	}
	// Apply the network adapter guestNetworkAdapterConfiguration
	virtualMachine, err := vna.GetVirtualMachine()
	if err != nil {
		return
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	err = vmms.SetGuestNetworkAdapterConfiguration(virtualMachine.computerSystem, guestNetworkAdapterConfiguration)
	return
}

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
