package storage

import (
	"testing"
)

func TestImageManagementService_GetVHDSnapshotInformation(t *testing.T) {
	ims, err := LocalImageManagementService()
	if err != nil {
		t.Fatalf("LocalImageManagementService failed: %v", err)
	}
	path := `D:\Hyper-V\Virtual Hard Disks\test_hyperv_vhd_5E46D170-5D14-4C40-A837-B8BDBF71C59B.avhdx`
	//hardDisk, err := disk.GetVirtualHardDiskByPath(ims.Session, `D:\Hyper-V\Virtual Hard Disks\test_hyperv_vhd_5E46D170-5D14-4C40-A837-B8BDBF71C59B.avhdx`)
	//if err != nil {
	//	t.Fatalf("GetVirtualHardDiskByPath failed: %v", err)
	//}
	virtualHardDiskSettingData, err := ims.GetVirtualHardDiskSettingData(path)
	if err != nil {
		t.Fatalf("GetVirtualHardDiskSettingData failed: %v", err)
	}
	t.Logf("Virtual hard disk setting data: %v", virtualHardDiskSettingData)
}
