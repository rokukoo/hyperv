package hypervctl

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var (
	ok                 bool
	virtualMachine     *VirtualMachine
	findVirtualMachine *VirtualMachine
	err                error
	cpuCoreCount       int
	memorySizeMB       int
	vmName             string = "hypervctl_test_hyperv_vm"
	force              bool
)

var hypervPath = `D:\Hyper-V\`

func TestHyperVIntegration(t *testing.T) {

}

func TestCreateVirtualMachine(t *testing.T) {
	t.Log("TestCreateVirtualMachine")
	savePath := hypervPath + vmName

	// TestCreateVM
	cpuCoreCount = 2
	memorySizeMB = 2048

	// Build a virtual machine with given name, save path, cpu core count and memory size
	if virtualMachine, err = CreateVirtualMachine(vmName, savePath, cpuCoreCount, memorySizeMB); err != nil {
		t.Fatalf("CreateVirtualMachine failed: %v", err)
	}
	// Check the virtual machine status
	assert.Equal(t, StateStopped, virtualMachine.State())
	t.Logf("Virtual machine created: %v", virtualMachine)
	// Check the physical virtual machine
	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	// Check the virtual machine status
	assert.NotNil(t, findVirtualMachine)
	assert.EqualExportedValues(t, findVirtualMachine, virtualMachine)
}

func TestStartVM(t *testing.T) {
	t.Log("TestStartVM")

	// TestStartVM
	if err = virtualMachine.Start(); err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	assert.Equal(t, StateRunning, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, StateRunning, findVirtualMachine.State())
	assert.EqualExportedValues(t, findVirtualMachine, virtualMachine)
}

func TestStopVM(t *testing.T) {
	t.Log("TestStopVM")

	// TestStopVM
	force = true
	if err = virtualMachine.Stop(force); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	assert.Equal(t, StateStopped, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, StateStopped, findVirtualMachine.State())
	assert.EqualExportedValues(t, findVirtualMachine, virtualMachine)
}

func TestRebootVM(t *testing.T) {
	t.Log("TestRebootVM")

	// TestRebootVM
	force = true
	if err = virtualMachine.Reboot(force); err != nil {
		t.Fatalf("Reboot failed: %v", err)
	}
	assert.Equal(t, StateRunning, virtualMachine.State())
}

func TestSuspendVM(t *testing.T) {
	t.Log("TestSuspendVM")

	// TestSuspendVM
	if err = virtualMachine.Suspend(); err != nil {
		t.Fatalf("Suspend failed: %v", err)
	}
	assert.Equal(t, StateSuspend, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, StateSuspend, findVirtualMachine.State())
}

func TestResumeVM(t *testing.T) {
	t.Log("TestResumeVM")

	// TestResumeVM
	if err = virtualMachine.Resume(); err != nil {
		t.Fatalf("Resume failed: %v", err)
	}
	assert.Equal(t, StateRunning, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, StateRunning, findVirtualMachine.State())
}

func TestModifyVM_cpuCoreCount(t *testing.T) {
	t.Log("TestModifyVM_cpuCoreCount")

	TestStartVM(t)

	// TestModifyVM cpuCoreCount
	cpuCoreCount = 4
	assert.Equal(t, StateRunning, virtualMachine.State())
	if ok, err = ModifyVirtualMachineSpecByName(vmName, cpuCoreCount, 0); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, StateRunning, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, cpuCoreCount, findVirtualMachine.CpuCoreCount)

	//assert.EqualExportedValues(t, virtualMachine, findVirtualMachine)
}

func TestModifyVM_memorySizeMB(t *testing.T) {
	t.Log("TestModifyVM_memorySizeMB")

	// TestModifyVM memorySizeMB
	memorySizeMB = 4096
	assert.Equal(t, StateRunning, virtualMachine.State())
	if ok, err = ModifyVirtualMachineSpecByName(vmName, 0, memorySizeMB); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, StateRunning, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, memorySizeMB, findVirtualMachine.MemorySizeMB)

	//assert.EqualExportedValues(t, virtualMachine, findVirtualMachine)
}

func TestModifyVM_cpuCoreCount_and_memorySizeMB(t *testing.T) {
	t.Log("TestModifyVM_cpuCoreCount_and_memorySizeMB")

	TestStartVM(t)

	// TestModifyVM cpuCoreCount and memorySizeMB
	cpuCoreCount = 2
	memorySizeMB = 2048
	assert.Equal(t, StateRunning, virtualMachine.State())
	if ok, err = ModifyVirtualMachineSpecByName(vmName, cpuCoreCount, memorySizeMB); err != nil {
		t.Fatalf("ModifyVM failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, StateRunning, virtualMachine.State())

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FindVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, cpuCoreCount, findVirtualMachine.CpuCoreCount)
	assert.Equal(t, memorySizeMB, findVirtualMachine.MemorySizeMB)

	//assert.EqualExportedValues(t, virtualMachine, findVirtualMachine)
}

func TestDeleteVM(t *testing.T) {
	t.Log("TestDeleteVM")

	if virtualMachine != nil {
		if err = virtualMachine.ForceStop(); err != nil {
			t.Fatalf("ForceStop failed: %v", err)
		}
	}

	// TestDeleteVM
	if ok, err = DeleteVirtualMachineByName(vmName); err != nil {
		t.Fatalf("DeleteVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, true, ok)
	t.Log("Virtual machine deleted")

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil && !errors.Is(err, wmiext.NotFound) {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	assert.Nil(t, findVirtualMachine)
}

func TestVMIntegration(t *testing.T) {
	t.Run("TestCreateVirtualMachine", TestCreateVirtualMachine)
	t.Run("TestStartVM", TestStartVM)
	t.Run("TestStopVM", TestStopVM)
	t.Run("TestRebootVM", TestRebootVM)
	t.Run("TestSuspendVM", TestSuspendVM)
	t.Run("TestResumeVM", TestResumeVM)
	t.Run("TestStopVM", TestStopVM)
	t.Run("TestModifyVM_cpuCoreCount", TestModifyVM_cpuCoreCount)
	t.Run("TestModifyVM_memorySizeMB", TestModifyVM_memorySizeMB)
	t.Run("TestModifyVM_cpuCoreCount_and_memorySizeMB", TestModifyVM_cpuCoreCount_and_memorySizeMB)
	t.Run("TestDeleteVM", TestDeleteVM)
}

func ExampleDeleteVirtualMachineByName() {
	if _, err = DeleteVirtualMachineByName(vmName); err != nil {
		log.Fatalf("DeleteVirtualMachineByName failed: %v", err)
	}
	fmt.Println("Virtual machine deleted")
	// Output:
	// Virtual machine deleted
}

func MustFindTestVirtualMachine() *VirtualMachine {
	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		log.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	return findVirtualMachine
}
