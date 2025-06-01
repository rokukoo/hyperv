package virtual_system

import (
	"testing"
)

//func TestGetHostComputerSystem(t *testing.T) {
//	host, err := GetHostComputerSystem()
//	if err != nil {
//		t.Errorf("Error getting host computer system: %v", err)
//		return
//	}
//	assert.NotNil(t, host)
//	t.Log(host)
//}

func TestVirtualSystemManagementService_ListComputerSystems(t *testing.T) {
	service, err := LocalVirtualSystemManagementService()
	if err != nil {
		t.Fatalf("Error getting virtual system management service: %v", err)
	}
	vms, err := service.ListComputerSystems()
	if err != nil {
		t.Fatalf("Error listing computer systems: %v", err)
	}
	for _, vm := range vms {
		t.Logf("VM: %v", vm)
	}
}

func TestComputerSystem_State(t *testing.T) {
	vmName := "hyperv_test_hyperv_vm"
	service, err := LocalVirtualSystemManagementService()
	if err != nil {
		t.Fatalf("Error getting virtual system management service: %v", err)
	}
	vm, err := service.FirstComputerSystemByName(vmName)
	if err != nil {
		t.Fatalf("Error getting VM by name: %v", err)
	}
	state, err := vm.GetState()
	if err != nil {
		t.Fatalf("Error getting state of VM: %v", err)
	}
	t.Logf("VM state: %v", state)

	//if err = vm.Start(); err != nil {
	//	t.Fatalf("Error getting state of VM: %v", err)
	//}
	//if err = vm.Stop(false); err != nil {
	//	t.Fatalf("Error getting state of VM: %v", err)
	//}
	//if err = vm.Start(); err != nil {
	//	t.Fatalf("Error getting state of VM: %v", err)
	//}
	//if err = vm.Reboot(false); err != nil {
	//	t.Fatalf("Error getting state of VM: %v", err)
	//}
}
