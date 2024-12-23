package wmiext

import (
	"testing"
)

func TestListPhysicalNetAdapter(t *testing.T) {
	t.Log("TestListPhysicalNetAdapter")
	physicalNetAdapters, err := ListPhysicalNetAdapter()
	if err != nil {
		t.Fatalf("ListPhysicalNetAdapter failed: %v", err)
	}
	for _, adapter := range physicalNetAdapters {
		//t.Logf("Physical network adapter: %v", adapter)
		t.Logf("Physical network adapter: Name=%s, InterfaceDescription=%s, InterfaceIndex=%d, Description=%s", adapter.Name, adapter.InterfaceDescription, adapter.InterfaceIndex, adapter.Description)
	}
	// Output:
	// Physical network adapter: Name=WLAN, InterfaceDescription=MediaTek Wi-Fi 6 MT7921 Wireless LAN Card, InterfaceIndex=7, Description=
	// Physical network adapter: Name=以太网, InterfaceDescription=Realtek PCIe GbE Family Controller, InterfaceIndex=5, Description=

}
