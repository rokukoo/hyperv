package hypervctl

import (
	"errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	virtualSwitchNames = struct {
		Private        string
		Internal       string
		ExternalBridge string
		ExternalDirect string
	}{
		Private:        "test_hyperv_vsw_private",
		Internal:       "test_hyperv_vsw_internal",
		ExternalBridge: "test_hyperv_vsw_bridge",
		ExternalDirect: "test_hyperv_vsw_external",
	}

	currentVirtualSwitchName string

	virtualSwitch        *VirtualSwitch
	findVirtualSwitch    *VirtualSwitch
	virtualSwitchType    VirtualSwitchType
	networkInterfaceName = "Realtek PCIe GbE Family Controller"
)

func TestCreatePrivateVirtualSwitch(t *testing.T) {
	// TestCreatePrivateVirtualSwitch
	t.Log("TestCreatePrivateVirtualSwitch")
	currentVirtualSwitchName = virtualSwitchNames.Private

	if virtualSwitch, err = CreateVirtualSwitch(currentVirtualSwitchName, VirtualSwitchTypePrivate, nil); err != nil {
		t.Fatalf("CreatePrivateVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Private virtual switch created: %v", virtualSwitch)
	}
	assert.NotNil(t, virtualSwitch)
	assert.EqualExportedValues(t, &VirtualSwitch{
		Name:            currentVirtualSwitchName,
		Description:     "",
		Type:            VirtualSwitchTypePrivate,
		PhysicalAdapter: nil,
	}, virtualSwitch)

	if findVirtualSwitch, err = FirstVirtualSwitchByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.EqualExportedValues(t, virtualSwitch, findVirtualSwitch)
}

func TestDeleteVirtualSwitchByName(t *testing.T) {
	// TestDeleteVirtualMachineByName
	t.Logf("TestDeleteVirtualMachineByName: name=%v", currentVirtualSwitchName)

	if err = DeleteVirtualSwitchByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("DeleteVirtualSwitchByName failed: %v", err)
	}

	if findVirtualSwitch, err = FirstVirtualSwitchByName(currentVirtualSwitchName); err != nil && !errors.Is(err, wmiext.NotFound) {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.Nil(t, findVirtualSwitch)
}

func TestChangeVirtualSwitchTypeByName(t *testing.T) {
	t.Log("TestChangeVirtualSwitchTypeByName")
	var originalType VirtualSwitchType

	if originalType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	t.Logf("Original virtual switch type: %v", originalType)

	t.Logf("ChangeVirtualSwitch [%v] from [%v] to [%v]", currentVirtualSwitchName, originalType, VirtualSwitchTypePrivate)
	if err = ChangeVirtualSwitchTypeByName(currentVirtualSwitchName, VirtualSwitchTypePrivate, nil); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypePrivate)

	t.Logf("ChangeVirtualSwitch [%v] from [%v] to [%v]", currentVirtualSwitchName, VirtualSwitchTypePrivate, VirtualSwitchTypeInternal)
	if err = ChangeVirtualSwitchTypeByName(currentVirtualSwitchName, VirtualSwitchTypeInternal, nil); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeInternal)

	t.Logf("ChangeVirtualSwitch [%v] from [%v] to [%v] with networkInterfaceName=%v", currentVirtualSwitchName, VirtualSwitchTypeInternal, VirtualSwitchTypeExternalBridge, networkInterfaceName)
	if err = ChangeVirtualSwitchTypeByName(currentVirtualSwitchName, VirtualSwitchTypeExternalBridge, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeExternalBridge)

	t.Logf("ChangeVirtualSwitch [%v] from [%v] to [%v] with networkInterfaceName=%v", currentVirtualSwitchName, VirtualSwitchTypeExternalBridge, VirtualSwitchTypeExternalDirect, networkInterfaceName)
	if err = ChangeVirtualSwitchTypeByName(currentVirtualSwitchName, VirtualSwitchTypeExternalDirect, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	virtualSwitchType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName)
	if err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, VirtualSwitchTypeExternalDirect)

	t.Logf("ChangeVirtualSwitch [%v] to original type [%v]", currentVirtualSwitchName, originalType)
	if err = ChangeVirtualSwitchTypeByName(currentVirtualSwitchName, originalType, &networkInterfaceName); err != nil {
		t.Fatalf("ChangeVirtualSwitchTypeByName failed: %v", err)
	}
	if virtualSwitchType, err = GetVirtualSwitchTypeByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("GetVirtualSwitchTypeByName failed: %v", err)
	}
	assert.Equal(t, virtualSwitchType, originalType)
}

func TestCreateInternalVirtualSwitch(t *testing.T) {
	// TestCreateInternalVirtualSwitch
	t.Log("TestCreateInternalVirtualSwitch")
	currentVirtualSwitchName = virtualSwitchNames.Internal

	if virtualSwitch, err = CreateVirtualSwitch(currentVirtualSwitchName, VirtualSwitchTypeInternal, nil); err != nil {
		t.Fatalf("CreateInternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("Internal virtual switch created: %v", virtualSwitch)
	}
	assert.NotNil(t, virtualSwitch)
	assert.EqualExportedValues(t, &VirtualSwitch{
		Name:            currentVirtualSwitchName,
		Description:     "",
		Type:            VirtualSwitchTypeInternal,
		PhysicalAdapter: nil,
	}, virtualSwitch)

	if findVirtualSwitch, err = FirstVirtualSwitchByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.EqualExportedValues(t, virtualSwitch, findVirtualSwitch)
}

func TestCreateBridgeVirtualSwitch(t *testing.T) {
	// TestCreateBridgeVirtualSwitch
	t.Log("TestCreateBridgeVirtualSwitch")
	currentVirtualSwitchName = virtualSwitchNames.ExternalBridge

	if virtualSwitch, err = CreateVirtualSwitch(currentVirtualSwitchName, VirtualSwitchTypeExternalBridge, &networkInterfaceName); err != nil {
		t.Fatalf("CreateBridgeVirtualSwitch failed: %v", err)
	} else {
		t.Logf("External bridge virtual switch created: %v", virtualSwitch)
	}
	assert.NotNil(t, virtualSwitch)
	assert.EqualExportedValues(t, &VirtualSwitch{
		Name:            currentVirtualSwitchName,
		Description:     "",
		Type:            VirtualSwitchTypeExternalBridge,
		PhysicalAdapter: &networkInterfaceName,
	}, virtualSwitch)

	if findVirtualSwitch, err = FirstVirtualSwitchByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.EqualExportedValues(t, virtualSwitch, findVirtualSwitch)
}

func TestCreateExternalVirtualSwitch(t *testing.T) {
	// TestCreateExternalVirtualSwitch
	t.Log("TestCreateExternalVirtualSwitch")
	currentVirtualSwitchName = virtualSwitchNames.ExternalDirect

	if virtualSwitch, err = CreateVirtualSwitch(currentVirtualSwitchName, VirtualSwitchTypeExternalDirect, &networkInterfaceName); err != nil {
		t.Fatalf("CreateExternalVirtualSwitch failed: %v", err)
	} else {
		t.Logf("External direct virtual switch created: %v", virtualSwitch)
	}
	assert.NotNil(t, virtualSwitch)
	assert.EqualExportedValues(t, &VirtualSwitch{
		Name:            currentVirtualSwitchName,
		Description:     "",
		Type:            VirtualSwitchTypeExternalDirect,
		PhysicalAdapter: &networkInterfaceName,
	}, virtualSwitch)

	if findVirtualSwitch, err = FirstVirtualSwitchByName(currentVirtualSwitchName); err != nil {
		t.Fatalf("FirstVirtualSwitchByName failed: %v", err)
	}
	assert.EqualExportedValues(t, virtualSwitch, findVirtualSwitch)
}

func TestPrivateVirtualSwitch(t *testing.T) {
	// TestPrivateVirtualSwitch
	t.Log("TestPrivateVirtualSwitch")
	t.Run("TestCreatePrivateVirtualSwitch", TestCreatePrivateVirtualSwitch)
	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)
	t.Run("TestDeleteVirtualSwitchByName", TestDeleteVirtualSwitchByName)
}

func TestInternalVirtualSwitch(t *testing.T) {
	// TestInternalVirtualSwitch
	t.Log("TestInternalVirtualSwitch")
	t.Run("TestCreateInternalVirtualSwitch", TestCreateInternalVirtualSwitch)
	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)
	t.Run("TestDeleteVirtualSwitchByName", TestDeleteVirtualSwitchByName)
}

func TestBridgeVirtualSwitch(t *testing.T) {
	t.Log("TestBridgeVirtualSwitch")
	t.Run("TestCreateBridgeVirtualSwitch", TestCreateBridgeVirtualSwitch)
	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)
	t.Run("TestDeleteVirtualSwitchByName", TestDeleteVirtualSwitchByName)
}

func TestExternalVirtualSwitch(t *testing.T) {
	// TestExternalVirtualSwitch
	t.Log("TestExternalVirtualSwitch")
	t.Run("TestCreateExternalVirtualSwitch", TestCreateExternalVirtualSwitch)
	t.Run("TestChangeVirtualSwitchTypeByName", TestChangeVirtualSwitchTypeByName)
	t.Run("TestDeleteVirtualSwitchByName", TestDeleteVirtualSwitchByName)
}

func TestVirtualSwitchIntegration(t *testing.T) {
	t.Run("TestPrivateVirtualSwitch", TestPrivateVirtualSwitch)
	t.Run("TestInternalVirtualSwitch", TestInternalVirtualSwitch)
	t.Run("TestBridgeVirtualSwitch", TestBridgeVirtualSwitch)
	t.Run("TestExternalVirtualSwitch", TestExternalVirtualSwitch)
}
