package hypervctl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	testVirtualHardDisk     *VirtualHardDisk
	testVirtualHardDiskName = "test_hyperv_vhd.vhdx"
	testVirtualHardDiskPath = hypervPath + `Virtual Hard Disks\` + testVirtualHardDiskName
	testVirtualHardDiskSize = 10 // 10GB
)

func TestVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	// TestCreateVirtualHardDisk
	t.Run("TestCreateVirtualHardDisk", TestCreateVirtualHardDisk)
	// TestDeleteVirtualHardDiskByPath
	t.Run("TestDeleteVirtualHardDiskByPath", TestDeleteVirtualHardDiskByPath)
}

func TestCreateVirtualHardDisk(t *testing.T) {
	t.Log("TestVirtualHardDisk")

	if testVirtualHardDisk, err = CreateVirtualHardDisk(testVirtualHardDiskPath, testVirtualHardDiskName, testVirtualHardDiskSize); err != nil {
		t.Fatalf("CreateVHD failed: %v", err)
	} else {
		t.Logf("VHD created: %v", testVirtualHardDisk)
	}

	assert.NotNil(t, testVirtualHardDiskPath)
	assert.Equal(t, true, checkVirtualHardDiskExistsByPath(testVirtualHardDiskPath))
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

func TestGetVirtualHardDiskSettingData(t *testing.T) {
	vhdxPath := `D:\\Hyper-V\\Virtual Hard Disks\\新建虚拟硬盘.vhdx`
	vhdSettingData, err := GetVirtualHardDiskSettingData(vhdxPath)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("vhdSettingData: %v", vhdSettingData)
}
