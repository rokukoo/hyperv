package hypervctl

import "testing"

func TestVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	// TestCreateVHD
	vhdName := "test_hyperv_vhd"
	vhdPath := hypervPath
	vhdSize := 10 // 10GB
	if vhd, err := CreateVirtualHardDisk(vhdName, vhdPath, vhdSize); err != nil {
		t.Fatalf("CreateVHD failed: %v", err)
	} else {
		t.Logf("VHD created: %v", vhd)
	}
	
}
