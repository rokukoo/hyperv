package wmictl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/virtualization/core/service"
	v2 "github.com/microsoft/wmi/server2019/root/virtualization/v2"
)

func GetVirtualSystemManagementService(wHost *host.WmiHost) (*service.VirtualSystemManagementService, error) {
	return newLocalVirtualSystemManagementService(wHost)
}

func newLocalVirtualSystemManagementService(whost *host.WmiHost) (mgmt *service.VirtualSystemManagementService, err error) {
	creds := whost.GetCredential()
	wQuery := query.NewWmiQuery("Msvm_VirtualSystemManagementService")
	// Refer to service.GetVirtualSystemManagementService exists an implicit bug, cache will not be eliminated if the service has been closed
	// So next time when you call GetVirtualSystemManagementService, it will return the closed service
	// This is a bug of the original code, so I turn to my own implementation
	// By the way, f**k microsoft! (╯▔皿▔)╯
	vmmswmi, err := v2.NewMsvm_VirtualSystemManagementServiceEx6(whost.HostName, string(constant.Virtualization), creds.UserName, creds.Password, creds.Domain, wQuery)
	if err != nil {
		return
	}
	mgmt = &service.VirtualSystemManagementService{Msvm_VirtualSystemManagementService: vmmswmi}
	return
}

func NewLocalVirtualSystemManagementService() (*service.VirtualSystemManagementService, error) {
	wmiHost := host.NewWmiLocalHost()
	vmms, err := newLocalVirtualSystemManagementService(wmiHost)
	if err != nil {
		return nil, err
	}
	return vmms, nil
}
