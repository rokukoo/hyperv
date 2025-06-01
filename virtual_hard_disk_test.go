package hyperv

import (
	"log"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestVirtualHardDiskValues struct {
	Name                string
	Size                float64
	VirtualHardDisk     *VirtualHardDisk
	FindVirtualHardDisk *VirtualHardDisk
}

func (tvhdv *TestVirtualHardDiskValues) Path() string {
	return hypervPath + `Virtual Hard Disks\` + tvhdv.Name
}

func NewTestVirtualHardDiskValues(size float64) *TestVirtualHardDiskValues {
	return &TestVirtualHardDiskValues{
		Name: uuid.NewString() + `.vhdx`,
		Size: size,
	}
}

var (
	testOsVirtualHardDisk   = NewTestVirtualHardDiskValues(10.0)
	testDataVirtualHardDisk = NewTestVirtualHardDiskValues(10.0)
)

var (
	currentVirtualHardDisk *TestVirtualHardDiskValues
)

func TestVirtualHardDiskIntegration(t *testing.T) {
	t.Log("TestVirtualHardDiskIntegration")
	t.Run("TestOsVirtualHardDisk", TestOsVirtualHardDisk)
	t.Run("TestDataVirtualHardDisk", TestDataVirtualHardDisk)
}

func TestVirtualMachine_GetVirtualHardDisks(t *testing.T) {
	t.Log("TestVirtualMachine_GetVirtualHardDisks")
	var virtualHardDisks []*VirtualHardDisk
	findVirtualMachine, err = FirstVirtualMachineByName("iECcequjDNcz1MTW")
	if err != nil {
		t.Fatalf("FirstVirtualMachineByName failed: %v", err)
	}
	virtualHardDisks, err = findVirtualMachine.GetVirtualHardDisks()
	if err != nil {
		t.Fatalf("GetVirtualHardDisks failed: %v", err)
	}
	for _, vhd := range virtualHardDisks {
		assert.NotNil(t, vhd)
		log.Printf("Attached Virtual Hard Disk: %s, Type: %d, SizeGB: %.2f/%.2f", vhd.Name, vhd.Type, vhd.UsedSizeGB, vhd.TotalSizeGB)
	}
}

func TestOsVirtualHardDisk(t *testing.T) {
	t.Log("TestOsVirtualHardDisk")
	currentVirtualHardDisk = testOsVirtualHardDisk
	// TestCreateVirtualHardDisks
	t.Run("TestCreateVirtualHardDisk", TestCreateVirtualHardDisk)
	// TestVirtualHardDisk_AttachAsSystemDisk
	t.Run("TestVirtualHardDisk_AttachAsSystemDisk", TestVirtualHardDisk_AttachAsSystemDisk)
	// TestDetachVirtualHardDisk
	t.Run("TestVirtualHardDisk_Detach", TestVirtualHardDisk_Detach)
	// TestVirtualHardDisk_Resize
	t.Run("TestVirtualHardDisk_Resize", TestVirtualHardDisk_Resize)
	//TestDeleteVirtualHardDiskByPath
	t.Run("TestDeleteVirtualHardDiskByPath", TestDeleteVirtualHardDiskByPath)
}

func TestDataVirtualHardDisk(t *testing.T) {
	t.Log("TestDataVirtualHardDisk")
	currentVirtualHardDisk = testDataVirtualHardDisk
	// TestCreateVirtualHardDisks
	t.Run("TestCreateVirtualHardDisk", TestCreateVirtualHardDisk)
	// TestVirtualHardDisk_AttachAsDataDisk
	t.Run("TestVirtualHardDisk_AttachAsDataDisk", TestVirtualHardDisk_AttachAsDataDisk)
	// TestVirtualHardDisk_Resize
	t.Run("TestVirtualHardDisk_Resize", TestVirtualHardDisk_Resize)
	// TestVirtualHardDisk_Resize_WithoutStop
	t.Run("TestVirtualHardDisk_Resize_WithoutStop", TestVirtualHardDisk_Resize_WithoutStop)
	// TestDetachVirtualHardDisk
	t.Run("TestVirtualHardDisk_Detach", TestVirtualHardDisk_Detach)
	//TestDeleteVirtualHardDiskByPath
	t.Run("TestDeleteVirtualHardDiskByPath", TestDeleteVirtualHardDiskByPath)
}

func TestCreateVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	require.NoDirExists(t, currentVirtualHardDisk.Path())
	assert.Equal(t, false, existsVirtualHardDiskByPath(currentVirtualHardDisk.Path()))

	currentVirtualHardDisk.VirtualHardDisk = nil
	currentVirtualHardDisk.FindVirtualHardDisk = nil

	if currentVirtualHardDisk.VirtualHardDisk, err = CreateVirtualHardDisk(currentVirtualHardDisk.Path(), currentVirtualHardDisk.Size); err != nil {
		t.Fatalf("CreateVHD failed: %v", err)
	} else {
		t.Logf("VHD created: %v", currentVirtualHardDisk.VirtualHardDisk)
	}

	if currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}

	assert.NotNil(t, currentVirtualHardDisk.FindVirtualHardDisk)
	assert.EqualExportedValues(t, currentVirtualHardDisk.VirtualHardDisk, currentVirtualHardDisk.FindVirtualHardDisk)
}

func TestVirtualHardDisk_AttachAsDataDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk_AttachAsDataDisk")
	currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path())
	findVirtualMachine = MustFindTestVirtualMachine(t)
	state := findVirtualMachine.State()
	if state != StateRunning {
		if err = findVirtualMachine.Start(); err != nil {
			t.Fatalf("StartVM failed: %v", err)
		}
	}
	defer func() {
		if state != StateRunning {
			if err = findVirtualMachine.computerSystem.ChangeState(state); err != nil {
				return
			}
		}
	}()
	t.Log("Test Attach Data Disk while VM is running")

	if ok, err = currentVirtualHardDisk.FindVirtualHardDisk.AttachAsDataDisk(findVirtualMachine); err != nil {
		t.Fatalf("AttachAsDataDisk failed: %v", err)
	}
	assert.Equal(t, true, ok)

	t.Run("TestVirtualHardDisk_Detach", TestVirtualHardDisk_Detach)

	t.Log("Force Stop VM and then Attach Data Disk")
	if err = findVirtualMachine.ForceStop(); err != nil {
		t.Fatalf("ForceStopVM failed: %v", err)
	}

	if ok, err = currentVirtualHardDisk.FindVirtualHardDisk.AttachAsDataDisk(findVirtualMachine); err != nil {
		t.Fatalf("AttachAsSystemDisk failed: %v", err)
	}

	assert.Equal(t, true, ok)
	assert.Equal(t, true, currentVirtualHardDisk.FindVirtualHardDisk.Attached)
}

func TestVirtualHardDisk_AttachAsSystemDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk_AttachAsSystemDisk")
	currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path())
	findVirtualMachine = MustFindTestVirtualMachine(t)
	if findVirtualMachine.State() != StateRunning {
		if err = findVirtualMachine.Start(); err != nil {
			t.Fatalf("StartVM failed: %v", err)
		}
	}
	t.Log("Test Attach System Disk while VM is running, it should be failed")

	ok, err = currentVirtualHardDisk.FindVirtualHardDisk.AttachAsSystemDisk(findVirtualMachine)
	assert.Equal(t, false, ok)
	assert.NotNil(t, err)

	t.Log("Force Stop VM and then Attach System Disk")
	if err = findVirtualMachine.ForceStop(); err != nil {
		t.Fatalf("ForceStopVM failed: %v", err)
	}

	if ok, err = currentVirtualHardDisk.FindVirtualHardDisk.AttachAsSystemDisk(findVirtualMachine); err != nil {
		t.Fatalf("AttachAsSystemDisk failed: %v", err)
	}

	assert.Equal(t, true, ok)
	assert.Equal(t, true, currentVirtualHardDisk.FindVirtualHardDisk.Attached)
}

func TestVirtualHardDisk_Detach(t *testing.T) {
	t.Log("TestVirtualHardDisk_Detach")

	currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path())
	if err = currentVirtualHardDisk.FindVirtualHardDisk.Detach(); err != nil {
		t.Fatalf("DetachVHD failed: %v", err)
	}

	assert.Equal(t, false, currentVirtualHardDisk.FindVirtualHardDisk.Attached)
}

func TestVirtualHardDisk_Resize(t *testing.T) {
	t.Log("TestVirtualHardDisk_Resize")

	if currentVirtualHardDisk.VirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, currentVirtualHardDisk.VirtualHardDisk)

	resizeGiB := float64(20)
	if ok, err = currentVirtualHardDisk.VirtualHardDisk.Resize(resizeGiB); err != nil {
		t.Fatalf("ResizeVHD failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, resizeGiB, currentVirtualHardDisk.VirtualHardDisk.TotalSizeGB)

	if currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, currentVirtualHardDisk.FindVirtualHardDisk)
	assert.EqualExportedValues(t, currentVirtualHardDisk.VirtualHardDisk, currentVirtualHardDisk.FindVirtualHardDisk)
}

func TestVirtualHardDisk_Resize_WithoutStop(t *testing.T) {
	t.Log("TestVirtualHardDisk_Resize_WithoutStop")

	if currentVirtualHardDisk.VirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, currentVirtualHardDisk.VirtualHardDisk)
	assert.Equal(t, true, currentVirtualHardDisk.VirtualHardDisk.Attached)

	findVirtualMachine = MustFindTestVirtualMachine(t)
	state := findVirtualMachine.State()
	if state != StateRunning {
		if err = findVirtualMachine.Start(); err != nil {
			t.Fatalf("StartVM failed: %v", err)
		}
	}
	defer func() {
		if state != StateRunning {
			if err = findVirtualMachine.computerSystem.ChangeState(state); err != nil {
				return
			}
		}
	}()
	t.Log("Test Resize VHD while VM is running")

	resizeGiB := float64(25)
	if ok, err = currentVirtualHardDisk.VirtualHardDisk.Resize(resizeGiB); err != nil {
		t.Fatalf("ResizeVHD failed: %v", err)
	}
	assert.Equal(t, true, ok)
	assert.Equal(t, resizeGiB, currentVirtualHardDisk.VirtualHardDisk.TotalSizeGB)

	if currentVirtualHardDisk.FindVirtualHardDisk, err = GetVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, currentVirtualHardDisk.FindVirtualHardDisk)
	assert.EqualExportedValues(t, currentVirtualHardDisk.VirtualHardDisk, currentVirtualHardDisk.FindVirtualHardDisk)
}

func TestDeleteVirtualHardDiskByPath(t *testing.T) {
	t.Log("TestDeleteVirtualHardDiskByPath")

	if ok, err = DeleteVirtualHardDiskByPath(currentVirtualHardDisk.Path()); err != nil {
		t.Fatalf("DeleteVHD failed: %v", err)
	} else {
		t.Logf("VHD deleted: %v", currentVirtualHardDisk.Path())
	}

	assert.Equal(t, true, ok)
	assert.Equal(t, false, existsVirtualHardDiskByPath(currentVirtualHardDisk.Path()))
}
