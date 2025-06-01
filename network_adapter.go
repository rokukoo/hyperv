package hyperv

import (
	"github.com/rokukoo/hyperv/pkg/hypervsdk/networking"
	utils "github.com/rokukoo/hyperv/pkg/hypervsdk/utils"
)

// ListAvailablePhysicalNetworkAdapters 列出所有可用的物理网络适配器
func ListAvailablePhysicalNetworkAdapters() ([]string, error) {
	var nics []string
	service, err := utils.NewLocalHyperVService()
	if err != nil {
		return nil, err
	}
	externalEthernetPorts, err := networking.ListEnabledExternalEthernetPort(service)
	if err != nil {
		return nil, err
	}
	for _, port := range externalEthernetPorts {
		nics = append(nics, port.Name)
	}

	wiFiPorts, err := networking.ListEnabledWiFiPort(service)
	if err != nil {
		return nil, err
	}
	for _, port := range wiFiPorts {
		nics = append(nics, port.Name)
	}

	return nics, nil
}
