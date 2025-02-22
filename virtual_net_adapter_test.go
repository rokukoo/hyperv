package hypervctl

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	//virtualNetworkAdapterName string                 = uuid.NewString()
	virtualNetworkAdapterName string                 = "50f9f144-0e4b-4827-979e-e0f4d4b7ff2a"
	virtualNetworkAdapter     *VirtualNetworkAdapter = &VirtualNetworkAdapter{
		Name:              virtualNetworkAdapterName,
		IsEnableBandwidth: true,
		MaxBandwidth:      10,
		MinBandwidth:      0,
	}
)

func TestVirtualNetworkAdapterIntegration(t *testing.T) {
	t.Log("Test virtualNetworkAdapter")
	t.Run("TestAddVirtualNetworkAdapter", TestVirtualMachine_AddVirtualNetworkAdapter)
	t.Run("TestConnect", TestVirtualNetworkAdapter_Connect)
	t.Run("TestDisConnect", TestVirtualNetworkAdapter_DisConnect)
	t.Run("TestRemoveVirtualNetworkAdapter", TestVirtualMachine_RemoveVirtualNetworkAdapter)
}

func TestVirtualMachine_AddVirtualNetworkAdapter(t *testing.T) {
	findVirtualMachine = MustFindTestVirtualMachine(t)
	if err = findVirtualMachine.AddVirtualNetworkAdapter(virtualNetworkAdapter); err != nil {
		t.Fatalf("AddVirtualNetworkAdapter failed: %v", err)
	}
	t.Logf("Virtual network adapter added successfully")
}

func TestVirtualMachine_RemoveVirtualNetworkAdapter(t *testing.T) {
	findVirtualMachine = MustFindTestVirtualMachine(t)
	if err = findVirtualMachine.RemoveVirtualNetworkAdapter(virtualNetworkAdapterName); err != nil {
		t.Fatalf("RemoveVirtualNetworkAdapter failed: %v", err)
	}
	t.Logf("Virtual network adapter removed successfully")
}

func TestVirtualNetworkAdapter_Connect(t *testing.T) {
	if virtualNetworkAdapter, err = FirstVirtualNetworkAdapterByName(virtualNetworkAdapterName); err != nil {
		t.Fatalf("FirstVirtualNetworkAdapterByName failed: %v", err)
	}
	if ok, err = virtualNetworkAdapter.ConnectByName(virtualSwitchNames.Private); err != nil {
		t.Fatalf("ConnectByName failed: %v", err)
	}
	t.Logf("Virtual network adapter connected successfully")
}

func TestVirtualNetworkAdapter_DisConnect(t *testing.T) {
	if virtualNetworkAdapter, err = FirstVirtualNetworkAdapterByName(virtualNetworkAdapterName); err != nil {
		t.Fatalf("FirstVirtualNetworkAdapterByName failed: %v", err)
	}
	if err = virtualNetworkAdapter.DisConnect(); err != nil {
		t.Fatalf("DisConnect failed: %v", err)
	}
	t.Logf("Virtual network adapter disconnected successfully")
}

//func TestVirtualNetworkAdapter_SetBandwidthOut(t *testing.T) {
//	if virtualNetworkAdapter, err = FirstVirtualNetworkAdapterByName(virtualNetworkAdapterName); err != nil {
//		t.Fatalf("FirstVirtualNetworkAdapterByName failed: %v", err)
//	}
//	err = virtualNetworkAdapter.SetBandwidthOut(100, 0)
//	if err != nil {
//		t.Fatalf("SetBandwidthOut failed: %v", err)
//	}
//}

func TestVirtualNetworkAdapter_SetBandwidthOut(t *testing.T) {
	if virtualNetworkAdapter, err = FirstVirtualNetworkAdapterByName(virtualNetworkAdapterName); err != nil {
		t.Fatalf("FirstVirtualNetworkAdapterByName failed: %v", err)
	}
	for _, limit := range []float64{-1, 0, 1} {
		for _, reserve := range []float64{-1, 0, 1} {
			actualLimit := limit
			actualReserve := reserve
			t.Logf("SetBandwidthOut limit=%v reserve=%v", limit, reserve)
			err = virtualNetworkAdapter.SetBandwidth(limit, reserve)
			if limit < 0 {
				actualLimit = 0
			}
			if reserve < 0 {
				actualReserve = 0
			}
			if err != nil {
				if errors.Is(err, wmiext.NotSupported) {
					continue
				}
				t.Fatalf("SetBandwidthOut failed: %v", err)
			}
			assert.Equal(t, actualLimit, virtualNetworkAdapter.MaxBandwidth)
			assert.Equal(t, actualReserve, virtualNetworkAdapter.MinBandwidth)
		}
	}
}
