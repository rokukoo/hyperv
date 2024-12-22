package wmictl

import (
	"testing"
)

func TestGetVirtualHardDiskState(t *testing.T) {
	vhdxPath := `D:\Hyper-V\Virtual Hard Disks\新建虚拟硬盘.vhdx`
	state, err := GetVirtualHardDiskState(vhdxPath)
	if err != nil {
		t.Fatalf("GetVirtualHardDiskState failed: %v", err)
	}
	t.Logf("Virtual hard disk state: %v", state)
}
