package hypervctl

import (
	"errors"
	"fmt"
	errors2 "github.com/microsoft/wmi/pkg/errors"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	ok             bool
	virtualMachine *HyperVVirtualMachine
	find           *HyperVVirtualMachine
	err            error
	cpuCoreCount   int
	memorySizeMB   int
	vmName         string
	force          bool
)

var hypervPath = `D:\Hyper-V\`

func TestCreateVirtualMachine(t *testing.T) {
	t.Log("TestCreateVirtualMachine")

	vmName = "hypervctl_test_hyperv_vm"

	savePath := hypervPath + vmName

	// TestCreateVM
	cpuCoreCount = 2
	memorySizeMB = 2048

	// Create a virtual machine with given name, save path, cpu core count and memory size
	if virtualMachine, err = CreateVirtualMachine(vmName, savePath, cpuCoreCount, memorySizeMB); err != nil {
		t.Fatalf("CreateVirtualMachine failed: %v", err)
	}
	// Check the virtual machine status
	assert.Equal(t, VMStatusStopped, virtualMachine.Status)
	t.Logf("Virtual machine created: %v", virtualMachine)
	// Check the physical virtual machine
	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	// Check the virtual machine status
	assert.NotNil(t, find)
	assert.ObjectsAreEqualValues(find, virtualMachine)
}

func TestStartVM(t *testing.T) {
	t.Log("TestStartVM")

	// TestStartVM
	if ok, err = virtualMachine.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)

	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, VMStatusRunning, find.Status)
	assert.ObjectsAreEqualValues(find, virtualMachine)
}

func TestStopVM(t *testing.T) {
	t.Log("TestStopVM")

	// TestStopVM
	force = true
	if ok, err = virtualMachine.Stop(force); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusStopped, virtualMachine.Status)

	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, VMStatusStopped, find.Status)
	assert.ObjectsAreEqualValues(find, virtualMachine)
}

func TestRebootVM(t *testing.T) {
	t.Log("TestRebootVM")

	TestStartVM(t)

	// TestRebootVM
	force = true
	if ok, err = virtualMachine.Reboot(force); err != nil {
		t.Fatalf("Reboot failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
}

func TestSuspendVM(t *testing.T) {
	t.Log("TestSuspendVM")

	// TestSuspendVM
	if ok, err = virtualMachine.Suspend(); err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusSuspended, virtualMachine.Status)

	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, VMStatusSuspended, find.Status)
}

func TestResumeVM(t *testing.T) {
	t.Log("TestResumeVM")

	// TestResumeVM
	if ok, err = virtualMachine.Resume(); err != nil {
		t.Fatalf("Resume failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)

	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, VMStatusRunning, find.Status)
}

func TestModifyVM_cpuCoreCount(t *testing.T) {
	t.Log("TestModifyVM_cpuCoreCount")

	// TestModifyVM cpuCoreCount
	cpuCoreCount = 4
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if ok, err = ModifyVirtualMachineSpec(vmName, cpuCoreCount, 0); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, cpuCoreCount, find.CpuCoreCount)
}

func TestModifyVM_memorySizeMB(t *testing.T) {
	t.Log("TestModifyVM_memorySizeMB")

	// TestModifyVM memorySizeMB
	memorySizeMB = 4096
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if ok, err = ModifyVirtualMachineSpec(vmName, 0, memorySizeMB); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, memorySizeMB, find.MemorySize)
}

func TestModifyVM_cpuCoreCount_and_memorySizeMB(t *testing.T) {
	t.Log("TestModifyVM_cpuCoreCount_and_memorySizeMB")

	// TestModifyVM cpuCoreCount and memorySizeMB
	cpuCoreCount = 2
	memorySizeMB = 2048
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if ok, err = ModifyVirtualMachineSpec(vmName, cpuCoreCount, memorySizeMB); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, VMStatusRunning, virtualMachine.Status)
	if find, err = GetVirtualMachineByName(vmName); err != nil {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, cpuCoreCount, find.CpuCoreCount)
}

func TestDeleteVM(t *testing.T) {
	// TestDeleteVM
	if ok, err = DeleteVirtualMachineByName(vmName, true); err != nil {
		t.Fatalf("DeleteVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
	t.Log("Virtual machine deleted")
	find, err = GetVirtualMachineByName(vmName)
	if err != nil && !errors.Is(err, errors2.NotFound) {
		t.Fatalf("GetVirtualMachineByName failed: %v", err)
	}
	assert.Nil(t, find)
}

func TestVMIntegration(t *testing.T) {
	t.Run("TestCreateVirtualMachine", TestCreateVirtualMachine)
	t.Run("TestStartVM", TestStartVM)
	t.Run("TestStopVM", TestStopVM)
	t.Run("TestRebootVM", TestRebootVM)
	t.Run("TestSuspendVM", TestSuspendVM)
	t.Run("TestResumeVM", TestResumeVM)
	t.Run("TestModifyVM_cpuCoreCount", TestModifyVM_cpuCoreCount)
	t.Run("TestModifyVM_memorySizeMB", TestModifyVM_memorySizeMB)
	t.Run("TestModifyVM_cpuCoreCount_and_memorySizeMB", TestModifyVM_cpuCoreCount_and_memorySizeMB)
	t.Run("TestDeleteVM", TestDeleteVM)
}

func ExampleDeleteVirtualMachineByName() {
	var (
		vmName string
		err    error
	)
	vmName = "test_hyperv_vm"
	if _, err = DeleteVirtualMachineByName(vmName, true); err != nil {
		log.Fatalf("DeleteVirtualMachineByName failed: %v", err)
	}
	fmt.Println("Virtual machine deleted")
	// Output:
	// Virtual machine deleted
}
