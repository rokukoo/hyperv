package hypervctl

import "testing"

var (
	virtualNetworkAdapter     *VirtualNetworkAdapter
	virtualNetworkAdapterName string = "Test"
)

func TestVirtualMachine_AddVirtualNetworkAdapter(t *testing.T) {
	findVirtualMachine = MustFindTestVirtualMachine()
	virtualNetworkAdapter = &VirtualNetworkAdapter{
		Name:              virtualNetworkAdapterName,
		IsEnableBandwidth: true,
		MaxBandwidth:      10,
	}
	if err = findVirtualMachine.AddVirtualNetworkAdapter(virtualNetworkAdapter); err != nil {
		t.Fatalf("AddVirtualNetworkAdapter failed: %v", err)
	}
	t.Logf("Virtual network adapter added successfully")
}

func TestVirtualMachine_RemoveVirtualNetworkAdapter(t *testing.T) {
	findVirtualMachine = MustFindTestVirtualMachine()
	if err = findVirtualMachine.RemoveVirtualNetworkAdapter(virtualNetworkAdapterName); err != nil {
		t.Fatalf("RemoveVirtualNetworkAdapter failed: %v", err)
	}
	t.Logf("Virtual network adapter removed successfully")
}

func TestVirtualNetworkAdapter_Connect(t *testing.T) {
	virtualNetworkAdapter = MustFirstVirtualNetworkAdapterByName(virtualNetworkAdapterName)
	vswName := "lan"
	if ok, err = virtualNetworkAdapter.ConnectByName(vswName); err != nil {
		t.Fatalf("ConnectByName failed: %v", err)
	}
	t.Logf("Virtual network adapter connected successfully")
}

func TestVirtualNetworkAdapter_DisConnect(t *testing.T) {
	virtualNetworkAdapter = MustFirstVirtualNetworkAdapterByName(virtualNetworkAdapterName)
	if err = virtualNetworkAdapter.DisConnect(); err != nil {
		t.Fatalf("DisConnect failed: %v", err)
	}
	t.Logf("Virtual network adapter disconnected successfully")
}
