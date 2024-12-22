package wmictl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
)

func GetVirtualMachineByVMName(vmName string) (*virtualsystem.VirtualMachine, error) {
	wHost := host.NewWmiLocalHost()
	return virtualsystem.GetVirtualMachineByVMName(wHost, vmName)
}
