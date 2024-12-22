package wmictl

import (
	"strconv"
	"testing"
)

func TestFindDataFileByPath(t *testing.T) {
	vhdxPath := `D:\Hyper-V\Virtual Hard Disks\新建虚拟硬盘.vhdx`
	dataFile, err := FindDataFileByPath(vhdxPath)
	if err != nil {
		t.Error(err)
	}
	fileSizeProperty, err := dataFile.GetProperty("FileSize")
	if err != nil {
		t.Error(err)
	}
	fileSize, err := strconv.Atoi(fileSizeProperty.(string))
	if err != nil {
		t.Error(err)
	}
	t.Logf("file size: %d", fileSize)
}
