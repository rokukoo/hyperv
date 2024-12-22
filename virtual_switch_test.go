package hypervctl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVSWIntegration(t *testing.T) {
	var (
		vsw  *VirtualSwitch
		find *VirtualSwitch
		ok   bool
		err  error
	)
	// TestCreatePrivateVirtualSwitch
	vswName := "test_hyperv_vsw"

	// Create a private virtual switch with given name
	privateVswName := vswName + "_private"
	if vsw, err = CreatePrivateVirtualSwitch(privateVswName); err != nil {
		t.Fatalf("CreatePrivateVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Private virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if find, err = FindVirtualSwitchByName(privateVswName); err != nil {
		t.Fatalf("FindVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, find)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(privateVswName); err != nil {
		t.Fatalf("RemoveVirtualSwitch failed: %v", err)
	}
	assert.Equal(t, true, ok)

	// TestCreateInternalVirtualSwitch
	internalVswName := vswName + "_internal"
	if vsw, err = CreateInternalVirtualSwitch(internalVswName); err != nil {
		t.Fatalf("CreateInternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Internal virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if find, err = FindVirtualSwitchByName(internalVswName); err != nil {
		t.Fatalf("FindVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, find)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(internalVswName); err != nil {
		t.Fatalf("RemoveVirtualSwitch failed: %v", err)
	}
	assert.Equal(t, true, ok)
}

func TestCreateExternalVirtualSwitch(t *testing.T) {
	var (
		vsw *VirtualSwitch
		err error
	)
	vswName := "test_hyperv_vsw_external"
	networkInterfaceName := "Realtek PCIe GbE Family Controller"
	if vsw, err = CreateExternalVirtualSwitch(vswName, networkInterfaceName, false); err != nil {
		t.Fatalf("CreateExternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("External virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	ok, err := DeleteVirtualSwitchByName(vswName)
	if err != nil {
		t.Fatalf("RemoveVirtualSwitch failed: %v", err)
	}
	assert.Equal(t, true, ok)
}
