package hyperv

import (
	"fmt"
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/rokukoo/hyperv/pkg/wmiext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	ok                 bool
	virtualMachine     *VirtualMachine
	findVirtualMachine *VirtualMachine
	err                error
	cpuCoreCount       int
	memorySizeMB       int
	vmName             string = "e54a2d62-a2a6-469b-998c-858d08929771"
	//vmName string = uuid.NewString()
	force bool
)

var hypervPath = `D:\Hyper-V\`

func TestHyperVIntegration(t *testing.T) {
	// Test VirtualSwitch
	t.Log("Test VirtualSwitch")
	t.Run("Test VirtualSwitch", TestVirtualSwitchIntegration)
	t.Run("TestCreatePrivateVirtualSwitch", TestCreatePrivateVirtualSwitch)
	defer t.Run("TestDeleteVirtualSwitchByName", TestDeleteVirtualSwitchByName)

	// Test VirtualMachine
	t.Log("Test VirtualMachine")
	t.Run("Test VirtualMachine", func(t *testing.T) {
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
	})
	defer t.Run("TestDeleteVM", TestDeleteVM)

	// Test VirtualHardDisk
	t.Log("Test VirtualHardDisk")
	t.Run("Test VirtualHardDisk", TestVirtualHardDiskIntegration)

	// Test VirtualNetworkAdapter
	t.Run("Test VirtualNetworkAdapter", TestVirtualNetworkAdapterIntegration)
}

func TestCreateVirtualMachine(t *testing.T) {
	t.Log("TestCreateVirtualMachine")
	savePath := hypervPath + vmName
	require.NoDirExists(t, savePath)

	// TestCreateVM
	cpuCoreCount = 2
	memorySizeMB = 2048

	// Build a virtual machine with given name, save path, cpu core count and memory size
	t.Logf("Creating virtual machine: name=%v, savePath=%v, cpuCoreCount=%v, memorySizeMB=%v", vmName, savePath, cpuCoreCount, memorySizeMB)
	if virtualMachine, err = CreateVirtualMachine(vmName, savePath, cpuCoreCount, memorySizeMB); err != nil {
		t.Fatalf("CreateVirtualMachine failed: %v", err)
	}
	// Check the virtual machine status
	assert.Equal(t, StateStopped, virtualMachine.State())
	assert.EqualExportedValues(t, &VirtualMachine{
		Name:         vmName,
		Description:  "",
		SavePath:     savePath,
		CpuCoreCount: cpuCoreCount,
		MemorySizeMB: memorySizeMB,
	}, virtualMachine)
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

	if virtualMachine.State() != StateRunning {
		if err = virtualMachine.Start(); err != nil {
			t.Fatalf("VM Start failed: %v", err)
		}
	}

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

	if virtualMachine.State() != StateRunning {
		if err = virtualMachine.Start(); err != nil {
			t.Fatalf("VM Start failed: %v", err)
		}
	}

	// TestRebootVM
	force = true
	if err = virtualMachine.Reboot(force); err != nil {
		t.Fatalf("Reboot failed: %v", err)
	}
	assert.Equal(t, StateRunning, virtualMachine.State())
}

func TestSuspendVM(t *testing.T) {
	t.Log("TestSuspendVM")

	if virtualMachine.State() != StateRunning {
		if err = virtualMachine.Start(); err != nil {
			t.Fatalf("VM Start failed: %v", err)
		}
	}

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

	// Manual update
	virtualMachine.CpuCoreCount = cpuCoreCount

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

	// Manual update
	virtualMachine.MemorySizeMB = memorySizeMB

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
	// Manual update
	virtualMachine.CpuCoreCount = cpuCoreCount
	virtualMachine.MemorySizeMB = memorySizeMB

	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FindVirtualMachineByName failed: %v", err)
	}
	assert.Equal(t, cpuCoreCount, findVirtualMachine.CpuCoreCount)
	assert.Equal(t, memorySizeMB, findVirtualMachine.MemorySizeMB)

	assert.EqualExportedValues(t, virtualMachine, findVirtualMachine)
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

func MustFindTestVirtualMachine(t *testing.T) *VirtualMachine {
	if findVirtualMachine, err = FirstVirtualMachineByName(vmName); err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	return findVirtualMachine
}
