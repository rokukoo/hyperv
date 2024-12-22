package wmictl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
)

func GetWmiInstanceFromPath(namespaceName WmiNamespace, instancePath string) (*wmi.WmiInstance, error) {
	wHost := host.NewWmiLocalHost()
	return instance.GetWmiInstanceFromPath(wHost, namespaceName, instancePath)
}
