package hypervctl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testVirtualHardDisk     *VirtualHardDisk
	testVirtualHardDiskName = "test_hyperv_vhd.vhdx"
	testVirtualHardDiskPath = hypervPath + `Virtual Hard Disks\` + testVirtualHardDiskName
	testVirtualHardDiskSize = 10.0 // 10GB

	findVirtualHardDisk *VirtualHardDisk
)

func TestVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	// TestCreateVirtualHardDisks
	t.Run("TestCreateVirtualHardDisk", TestCreateVirtualHardDisk)
	// TestVirtualHardDisk_AttachTo
	t.Run("TestVirtualHardDisk_AttachTo", TestVirtualHardDisk_AttachTo)
	// TestDetachVirtualHardDisk
	t.Run("TestVirtualHardDisk_Detach", TestVirtualHardDisk_Detach)
	// TestVirtualHardDisk_AttachAsSystemDisk
	t.Run("TestVirtualHardDisk_AttachAsSystemDisk", TestVirtualHardDisk_AttachAsSystemDisk)
	// TestDetachVirtualHardDisk
	t.Run("TestVirtualHardDisk_Detach", TestVirtualHardDisk_Detach)
	// TestVirtualHardDisk_Resize
	//t.Run("TestVirtualHardDisk_Resize", TestVirtualHardDisk_Resize)
	//TestDeleteVirtualHardDiskByPath
	t.Run("TestDeleteVirtualHardDiskByPath", TestDeleteVirtualHardDiskByPath)
}

func TestVirtualHardDisk_AttachAsSystemDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk_AttachAsSystemDisk")

	findVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath)
	findVirtualMachine = MustFirstVirtualMachineByName(vmName)
	ok, err = findVirtualHardDisk.AttachAsSystemDisk(findVirtualMachine)
	if err != nil {
		t.Fatalf("AttachAsSystemDisk failed: %v", err)
	}
}

func TestVirtualHardDisk_Detach(t *testing.T) {
	t.Log("TestVirtualHardDisk_Detach")

	findVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath)
	if err = findVirtualHardDisk.Detach(); err != nil {
		t.Fatalf("DetachVHD failed: %v", err)
	}

	assert.Equal(t, false, findVirtualHardDisk.Attached)
}

func TestVirtualHardDisk_AttachTo(t *testing.T) {
	t.Log("TestVirtualHardDisk_AttachTo")

	findVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath)
	findVirtualMachine = MustFirstVirtualMachineByName(vmName)
	ok, err = findVirtualHardDisk.AttachTo(findVirtualMachine)
	if err != nil {
		t.Fatalf("AttachTo failed: %v", err)
		return
	}
}

func TestVirtualHardDisk_Resize(t *testing.T) {
	t.Log("TestVirtualHardDisk_Resize")

	if testVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, testVirtualHardDisk)

	resizeGiB := float64(20)
	if ok, err = testVirtualHardDisk.Resize(resizeGiB); err != nil {
		t.Fatalf("ResizeVHD failed: %v", err)
	}
	assert.Equal(t, true, ok)

	if findVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}
	assert.NotNil(t, findVirtualHardDisk)
	assert.EqualExportedValues(t, testVirtualHardDisk, findVirtualHardDisk)
}

func TestCreateVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	if checkVirtualHardDiskExistsByPath(testVirtualHardDiskPath) {
		t.Logf("test vhd exists: %v, try to delete it", testVirtualHardDiskPath)
		ok, err = DeleteVirtualHardDiskByPath(testVirtualHardDiskPath)
		if err != nil {
			t.Fatalf("DeleteVHD failed: %v", err)
		}

		assert.Equal(t, true, ok)
		assert.Equal(t, false, checkVirtualHardDiskExistsByPath(testVirtualHardDiskPath))
	}

	if testVirtualHardDisk, err = CreateVirtualHardDisk(testVirtualHardDiskPath, testVirtualHardDiskSize); err != nil {
		t.Fatalf("CreateVHD failed: %v", err)
	} else {
		t.Logf("VHD created: %v", testVirtualHardDisk)
	}

	if findVirtualHardDisk, err = GetVirtualHardDiskByPath(testVirtualHardDiskPath); err != nil {
		t.Fatalf("GetVHD failed: %v", err)
	}

	assert.NotNil(t, findVirtualHardDisk)
	assert.EqualExportedValues(t, testVirtualHardDisk, findVirtualHardDisk)
}

func TestDeleteVirtualHardDiskByPath(t *testing.T) {
	t.Log("TestDeleteVirtualHardDiskByPath")

	if ok, err = DeleteVirtualHardDiskByPath(testVirtualHardDiskPath); err != nil {
		t.Fatalf("DeleteVHD failed: %v", err)
	} else {
		t.Logf("VHD deleted: %v", testVirtualHardDiskPath)
	}

	assert.Equal(t, true, ok)
	assert.Equal(t, false, checkVirtualHardDiskExistsByPath(testVirtualHardDiskPath))
}

//func TestGetVirtualHardDiskSettingData(t *testing.T) {
//	vhdxPath := `D:\\Hyper-V\\Virtual Hard Disks\\新建虚拟硬盘.vhdx`
//	vhdSettingData, err := NewVirtualHardDiskSettingData(vhdxPath)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	t.Logf("vhdSettingData: %v", vhdSettingData)
//}
