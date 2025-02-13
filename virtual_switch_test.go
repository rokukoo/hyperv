package hypervctl

import (
	hypervsdk "github.com/rokukoo/hypervctl/pkg/hypervsdk/networking"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	vswName              string
	vsw                  *hypervsdk.VirtualEthernetSwitch
	findVsw              *hypervsdk.VirtualEthernetSwitch
	virtualSwitchType    VirtualSwitchType
	networkInterfaceName string = "Realtek PCIe GbE Family Controller"
)

func TestChangeVirtualSwitchTypeByName(t *testing.T) {
	var originalType VirtualSwitchType

	if originalType, err = GetVirtualSwitchTypeByName(vswName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}

	if err = ChangeVirtualSwitchTypeByName(vswName, VirtualSwitchTypePrivate, nil); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(vswName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypePrivate)

	if err = ChangeVirtualSwitchTypeByName(vswName, VirtualSwitchTypeInternal, nil); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(vswName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeInternal)

	if err = ChangeVirtualSwitchTypeByName(vswName, VirtualSwitchTypeExternalBridge, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(vswName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeExternalBridge)

	if err = ChangeVirtualSwitchTypeByName(vswName, VirtualSwitchTypeExternalDirect, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	virtualSwitchType, err = GetVirtualSwitchTypeByName(vswName)
	if err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeExternalDirect)

	if err = ChangeVirtualSwitchTypeByName(vswName, originalType, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(vswName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, originalType)
}

func TestPrivateVirtualSwitch(t *testing.T) {
	// TestCreatePrivateVirtualSwitch
	vswName = "test_hyperv_vsw_private"

	if vsw, err = CreateVirtualSwitch(vswName, "", VirtualSwitchTypePrivate, nil); err != nil {
		t.Fatalf("CreatePrivateVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Private virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if findVsw, err = FindVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, findVsw)

	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("DeleteVirtualSwitchByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
}

func TestInternalVirtualSwitch(t *testing.T) {
	// TestCreatePrivateVirtualSwitch
	vswName = "test_hyperv_vsw_internal"

	if vsw, err = CreateVirtualSwitch(vswName, "", VirtualSwitchTypeInternal, nil); err != nil {
		t.Fatalf("CreateInternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Internal virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if findVsw, err = FindVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, findVsw)

	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("DeleteVirtualSwitchByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
}

func TestBridgeVirtualSwitch(t *testing.T) {
	// TestCreatePrivateVirtualSwitch
	vswName = "test_hyperv_vsw_bridge"

	if vsw, err = CreateVirtualSwitch(vswName, "", VirtualSwitchTypeExternalBridge, &networkInterfaceName); err != nil {
		t.Fatalf("CreateBridgelVirtualSwitch failed: %v", err)
	} else {
		t.Logf("External bridge virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if findVsw, err = FindVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, findVsw)

	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("DeleteVirtualSwitchByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
}

func TestExternalVirtualSwitch(t *testing.T) {
	// TestCreatePrivateVirtualSwitch
	vswName = "test_hyperv_vsw_external"

	if vsw, err = CreateVirtualSwitch(vswName, "", VirtualSwitchTypeExternalDirect, &networkInterfaceName); err != nil {
		t.Fatalf("CreateExternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("External direct virtual switch created: %v", vsw)
	}
	assert.NotNil(t, vsw)

	if findVsw, err = FindVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.ObjectsAreEqualValues(vsw, findVsw)

	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)

	// TestDeleteVirtualMachineByName
	if ok, err = DeleteVirtualSwitchByName(vswName); err != nil {
		t.Fatalf("DeleteVirtualSwitchByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
}

func TestVSWIntegration(t *testing.T) {
	t.Run("TestPrivateVirtualSwitch", TestPrivateVirtualSwitch)
	t.Run("TestInternalVirtualSwitch", TestInternalVirtualSwitch)
	t.Run("TestBridgeVirtualSwitch", TestBridgeVirtualSwitch)
	t.Run("TestExternalVirtualSwitch", TestExternalVirtualSwitch)
}
