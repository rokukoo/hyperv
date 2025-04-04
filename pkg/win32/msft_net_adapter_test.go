package win32

import "testing"

func TestListNetworkAdapters(t *testing.T) {
	adapters, err := ListNetworkAdapters()
	if err != nil {
		t.Fatalf("failed to list network adapters: %v", err)
	}

	for _, adapter := range adapters {
		t.Logf("Adapter: %+v", adapter.Name)
	}
}

func TestGetNetworkAdapterByName(t *testing.T) {
	name := "vEthernet (lan)"
	adapter, err := GetNetworkAdapterByName(name)
	if err != nil {
		t.Fatalf("failed to get network adapter: %v", err)
	}

	t.Logf("Adapter: %+v", adapter)
}

func TestNetworkAdapter_Configure(t *testing.T) {
	name := "vEthernet (lan)"
	adapter, err := GetNetworkAdapterByName(name)
	if err != nil {
		t.Fatalf("failed to get network adapter: %v", err)
	}

	t.Logf("Adapter: %+v", adapter)

	ipAddress := []string{"172.16.0.2"}
	subnetMask := []string{"255.255.0.0"}
	gateway := []string{"172.16.0.1"}
	dnsServer := []string{"114.114.114.114", "8.8.8.8"}

	if err = adapter.Configure(
		ipAddress,
		subnetMask,
		gateway,
		dnsServer,
	); err != nil {
		t.Fatalf("failed to configure network adapter: %v", err)
	}

}
