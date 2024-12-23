package wmiext

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/base/session"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/pkg/hardware/network/netadapter"
	"github.com/microsoft/wmi/server2019/root/cimv2"
	"github.com/microsoft/wmi/server2019/root/standardcimv2"
	"github.com/pkg/errors"
	"strconv"
)

type NetworkAdapter struct {
	Name                 string
	InterfaceIndex       uint32
	Description          string
	MediaConnectState    string
	State                string
	InterfaceDescription string
	*netadapter.NetworkAdapter
}

func (nad *NetworkAdapter) Configure(
	ipaddress []string,
	subnetMask []string,
	gateway []string,
	dns []string,
) error {
	wmiHost := host.NewWmiLocalHost()
	hostSession, err := session.GetHostSession("root\\cimv2", wmiHost)
	if err != nil {
		return err
	}
	wquery := query.NewWmiQuery("Win32_NetworkAdapterConfiguration")
	wquery.AddFilter("InterfaceIndex", strconv.Itoa(int(nad.InterfaceIndex)))
	win32NetworkAdapterConfiguration, err := cimv2.NewWin32_NetworkAdapterConfigurationEx6(wmiHost.HostName, string(constant.CimV2), hostSession.Username, hostSession.Password, hostSession.Domain, wquery)
	if err != nil {
		return err
	}
	success, err := win32NetworkAdapterConfiguration.EnableStatic(ipaddress, subnetMask)
	if success != 0 {
		return errors.New("failed to enable static ip")
	}
	if err != nil {
		return err
	}
	// 设置 网关
	// NOTE: fucking golang and wmi, the real type for wmi 'uint16' in golang actually is 'uint8'
	gatewayCostMetrics := []uint8{1} // 网关跃点数
	retVal, err := win32NetworkAdapterConfiguration.InvokeMethodWithReturn("SetGateways", gateway, gatewayCostMetrics)
	if err != nil {
		return err
	}
	success = uint32(retVal)
	//success, err = win32NetworkAdapterConfiguration.SetGateways(gateway, gatewayCostMetrics)
	if success != 0 {
		return errors.New("failed to set gateway")
	}
	// 设置 DNS
	success, err = win32NetworkAdapterConfiguration.SetDNSServerSearchOrder(dns)
	if success != 0 {
		return errors.New("failed to set dns")
	}
	if err != nil {
		return err
	}
	return nil
}

func networkAdapter(nad *netadapter.NetworkAdapter) (*NetworkAdapter, error) {
	name, err := nad.GetProperty("Name")
	if err != nil {
		return nil, err
	}
	interfaceIndex, err := nad.GetInterfaceIndex()
	if err != nil {
		return nil, err
	}
	description, err := nad.GetPropertyDescription()
	if err != nil {
		return nil, err
	}

	// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/legacy/hh968170%28v=vs.85%29?redirectedfrom=MSDN
	mediaConnectState, err := nad.GetProperty("MediaConnectState")
	mediaConnectStateNum := mediaConnectState.(int32)
	if err != nil {
		return nil, err
	}
	var translatedMediaConnectState = "未知"
	if mediaConnectStateNum == 0 {
		translatedMediaConnectState = "未知"
	} else if mediaConnectStateNum == 1 {
		translatedMediaConnectState = "已连接"
	} else if mediaConnectStateNum == 2 {
		translatedMediaConnectState = "已断开"
	}

	state, err := nad.GetProperty("State")
	stateNum := state.(int32)
	if err != nil {
		return nil, err
	}
	var translatedState = "未知"
	if stateNum == 0 {
		translatedState = "未知"
	} else if stateNum == 1 {
		translatedState = "已存在"
	} else if stateNum == 2 {
		translatedState = "已启用"
	} else if stateNum == 3 {
		translatedState = "已禁用"
	}

	interfaceDescription, err := nad.GetPropertyInterfaceDescription()
	if err != nil {
		return nil, err
	}

	return &NetworkAdapter{
		Name:                 name.(string),
		InterfaceDescription: interfaceDescription,
		InterfaceIndex:       uint32(interfaceIndex),
		Description:          description,
		MediaConnectState:    translatedMediaConnectState,
		State:                translatedState,
		NetworkAdapter:       nad,
	}, nil
}

func FindNetAdapterByInterfaceDescription(interfaceDescription string) (*NetworkAdapter, error) {
	whost := host.NewWmiLocalHost()
	creds := whost.GetCredential()
	querytmp := query.NewWmiQuery("MSFT_NetAdapter", "InterfaceDescription", interfaceDescription)
	tmp, err := standardcimv2.NewMSFT_NetAdapterEx6(whost.HostName, string(constant.StadardCimV2), creds.UserName, creds.Password, creds.Domain, querytmp)
	if err != nil {
		return nil, err
	}
	nad := &netadapter.NetworkAdapter{MSFT_NetAdapter: tmp}
	return networkAdapter(nad)
}

func FindNetAdaptersByName(name string) ([]*NetworkAdapter, error) {
	var adapters []*NetworkAdapter
	whost := host.NewWmiLocalHost()
	wquery := query.NewWmiQuery("MSFT_NetAdapter")
	wquery.AddFilter("Name", name)
	collections, err := instance.GetWmiInstancesFromHost(whost, string(constant.StadardCimV2), wquery)
	if err != nil {
		return nil, err
	}
	adapters = make([]*NetworkAdapter, 0, len(collections))
	for _, nad := range collections {
		adapter, err := netadapter.NewNetworkAdapter(nad)
		if err != nil {
			return nil, err
		}
		value, err := adapter.GetInterfaceIndex()
		if err != nil {
			return nil, err
		}
		adapter.InterfaceIndex = uint32(value)
		netAdapter, err := networkAdapter(adapter)
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, netAdapter)
	}
	return adapters, nil
}

func ListPhysicalNetAdapter() ([]*NetworkAdapter, error) {
	var adapters []*NetworkAdapter
	whost := host.NewWmiLocalHost()
	wquery := query.NewWmiQuery("MSFT_NetAdapter")
	wquery.AddFilter("Virtual", "False")
	collections, err := instance.GetWmiInstancesFromHost(whost, string(constant.StadardCimV2), wquery)
	if err != nil {
		return nil, err
	}
	adapters = make([]*NetworkAdapter, 0, len(collections))
	for _, nad := range collections {
		adapter, err := netadapter.NewNetworkAdapter(nad)
		if err != nil {
			return nil, err
		}
		value, err := adapter.GetInterfaceIndex()
		if err != nil {
			return nil, err
		}
		adapter.InterfaceIndex = uint32(value)
		netAdapter, err := networkAdapter(adapter)
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, netAdapter)
	}
	return adapters, nil
}

func ListNetAdapters() ([]*NetworkAdapter, error) {
	var adapters []*NetworkAdapter
	whost := host.NewWmiLocalHost()
	wquery := query.NewWmiQuery("MSFT_NetAdapter")
	collections, err := instance.GetWmiInstancesFromHost(whost, string(constant.StadardCimV2), wquery)
	if err != nil {
		return nil, err
	}
	adapters = make([]*NetworkAdapter, 0, len(collections))
	for _, nad := range collections {
		adapter, err := netadapter.NewNetworkAdapter(nad)
		if err != nil {
			return nil, err
		}
		value, err := adapter.GetInterfaceIndex()
		if err != nil {
			return nil, err
		}
		adapter.InterfaceIndex = uint32(value)
		netAdapter, err := networkAdapter(adapter)
		if err != nil {
			return nil, err
		}
		adapters = append(adapters, netAdapter)
	}
	return adapters, nil
}
