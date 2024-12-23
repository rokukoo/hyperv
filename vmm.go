package hypervctl

import (
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	HyperVNamespace                = "root\\virtualization\\v2"
	VirtualSystemManagementService = "Msvm_VirtualSystemManagementService"
	MsvmComputerSystem             = "Msvm_ComputerSystem"
)

func NewLocalHyperVService() (*wmiext.Service, error) {
	service, err := wmiext.NewLocalService(HyperVNamespace)
	if err != nil {
		return nil, translateCommonHyperVWmiError(err)
	}

	return service, nil
}
