package hypervctl

import (
	errors2 "github.com/microsoft/wmi/pkg/errors"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	na "github.com/microsoft/wmi/pkg/virtualization/network/virtualnetworkadapter"
	"github.com/rokukoo/hypervctl/wmictl"
)

// AddNetworkAdapter adds a network adapter to the virtual machine
func (vm *HyperVVirtualMachine) AddNetworkAdapter(name string, limit int, reserve int) (*VirtualNetworkAdapter, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return nil, err
	}
	syntheticNetworkAdapter, err := virtualMachine.NewSyntheticNetworkAdapter(name)
	if err != nil {
		return nil, err
	}
	//service, err := wmictl.GetVirtualSystemManagementService(virtualMachine.GetWmiHost())
	service, err := wmictl.NewLocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}
	vmSettingData, err := virtualMachine.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	resultCollection, err := service.AddVirtualSystemResource(vmSettingData, syntheticNetworkAdapter.CIM_ResourceAllocationSettingData, -1)
	if err != nil {
		return nil, err
	}
	instance := resultCollection[0]
	adapter, err := na.NewVirtualNetworkAdapter(instance)
	if err != nil {
		return nil, err
	}
	vna, err := virtualNetworkAdapter(adapter)
	if err != nil {
		return nil, err
	}
	vna.VirtualMachine = vm

	if limit != 0 || reserve != 0 {
		if err = vna.SetBandwidth(float64(limit), float64(reserve)); err != nil {
			return nil, err
		}
	}
	return vna, nil
}

// GetNetworkAdapterByName returns a network interface by name
func (vm *HyperVVirtualMachine) GetNetworkAdapterByName(name string) (*VirtualNetworkAdapter, error) {
	virtualMachine, err := vm.VM()
	if err != nil {
		return nil, err
	}
	adapter, err := virtualMachine.GetVirtualNetworkAdapterByName(name)
	if err != nil {
		return nil, err
	}
	vNic, err := virtualNetworkAdapter(adapter)
	if err != nil {
		return nil, err
	}
	vNic.VirtualMachine = vm
	return vNic, nil
}

// RemoveNetworkAdapter removes a network interface by name
func (vm *HyperVVirtualMachine) RemoveNetworkAdapter(name string) (bool, error) {
	var (
		vmms           *service.VirtualSystemManagementService
		err            error
		virtualMachine *virtualsystem.VirtualMachine
		adapter        *na.VirtualNetworkAdapter
	)
	// Get the hyper-v physical virtual machine
	if virtualMachine, err = vm.VM(); err != nil {
		return false, err
	}
	// Get the virtual system management service
	if vmms, err = wmictl.NewLocalVirtualSystemManagementService(); err != nil {
		return false, err
	}
	// Get the network adapter by name
	if adapter, err = virtualMachine.GetVirtualNetworkAdapterByName(name); err != nil {
		return false, err
	}
	// If the adapter is not found, return an error
	if adapter == nil {
		return false, errors2.NotFound
	}
	// Remove the network adapter
	if err = vmms.RemoveVirtualNetworkAdapter(adapter); err != nil {
		return false, err
	}
	return true, nil
}
