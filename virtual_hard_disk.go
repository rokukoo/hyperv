package hypervctl

import (
	"fmt"
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/instance"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	wmierrors "github.com/microsoft/wmi/pkg/errors"
	"github.com/microsoft/wmi/pkg/virtualization/core/storage/disk"
	diskService "github.com/microsoft/wmi/pkg/virtualization/core/storage/service"
	"github.com/microsoft/wmi/pkg/virtualization/core/virtualsystem"
	wmi "github.com/microsoft/wmi/pkg/wmiinstance"
	v2 "github.com/microsoft/wmi/server2019/root/virtualization/v2"
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/wmictl"
	"os"
	"path/filepath"
	"slices"
)

// VirtualHardDisk
// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-virtualharddisksettingdata
type VirtualHardDisk struct {
	Name        string `json:"name"`
	TotalSizeGB uint64 `json:"total_size"`
	UsedSizeGB  uint64 `json:"used_size"`
	Path        string `json:"path"`
	*disk.VirtualHardDisk
}

func (vhd *VirtualHardDisk) Resize(newSizeGiB int) (ok bool, err error) {
	mgmt, err := diskService.GetImageManagementService(host.NewWmiLocalHost())
	if err != nil {
		return false, err
	}
	err = mgmt.ResizeDisk(vhd.Path, uint64(newSizeGiB)*1024*1024*1024)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (vhd *VirtualHardDisk) AttachTo(vm *HyperVVirtualMachine) (ok bool, err error) {
	return vhd.AttachToSCSI(vm)
}

func (vhd *VirtualHardDisk) AttachToIDE(vm *HyperVVirtualMachine) (ok bool, err error) {
	vmms, err := wmictl.NewLocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	virtualMachine, err := vm.VM()
	if err != nil {
		return
	}
	_, _, err = vmms.AttachVirtualHardDisk(virtualMachine, vhd.Path, virtualsystem.VirtualHardDiskType_OS_VIRTUALHARDDISK)
	if err != nil {
		return
	}
	return true, nil
}

func (vhd *VirtualHardDisk) AttachToSCSI(vm *HyperVVirtualMachine) (ok bool, err error) {
	vmms, err := wmictl.NewLocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	virtualMachine, err := vm.VM()
	if err != nil {
		return
	}
	if controllers, err := vm.GetSCSIControllers(); errors.Is(err, wmierrors.NotFound) || (controllers == nil || len(controllers) == 0) {
		if err = vmms.AddSCSIController(virtualMachine); err != nil {
			return false, err
		}
	} else if err != nil {
		return false, err
	}
	_, _, err = vmms.AttachVirtualHardDisk(virtualMachine, vhd.Path, virtualsystem.VirtualHardDiskType_DATADISK_VIRTUALHARDDISK)
	if err != nil {
		return
	}
	return true, nil
}

func (vhd *VirtualHardDisk) Detach() (ok bool, err error) {
	path := vhd.Path
	if !checkVirtualHardDiskExistsByPath(path) {
		return false, errors.New("vhd not exists")
	}
	wquery := query.NewWmiQuery("Msvm_StorageAllocationSettingData")
	// wquery.AddFilter("HostResource[0]", path)
	wquery.AddFilter("ResourceType", fmt.Sprintf("%d", v2.ResourcePool_ResourceType_Logical_Disk))
	wHost := host.NewWmiLocalHost()
	vhdSettings, err := instance.GetWmiInstancesFromHost(wHost, string(constant.Virtualization), wquery)
	var vhdSetting *wmi.WmiInstance
	for _, item := range vhdSettings {
		virtualHardDisk, err := disk.NewVirtualHardDisk(item)
		if err != nil {
			return false, err
		}
		hostResource, err := virtualHardDisk.GetPropertyHostResource()
		if err != nil {
			return false, err
		}
		if slices.Contains(hostResource, path) {
			vhdSetting = item
			break
		}
	}
	if vhdSetting == nil {
		return false, errors.New("vhd not mounted yet")
	}
	virtualHardDisk, err := disk.NewVirtualHardDisk(vhdSetting)
	vmms, err := wmictl.NewLocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	err = vmms.DetachVirtualHardDisk(virtualHardDisk)
	if err != nil {
		return
	}
	return true, nil
}

func GetVirtualHardDiskByPath(path string) (*VirtualHardDisk, error) {
	if !checkVirtualHardDiskExistsByPath(path) {
		return nil, errors.New("vhd not exists")
	}
	// 根据路径读取文件名
	fileName := filepath.Base(path)
	// 获取文件信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	// 获取文件大小 (以 GB 为单位)
	usedSizeGB := uint64(fileInfo.Size()) / (1024 * 1024 * 1024)

	// TODO: 通过 WMI 获取 VirtualHardDisk 信息, 拿到 TotalSizeGB

	return &VirtualHardDisk{
		Name:       fileName,
		UsedSizeGB: usedSizeGB,
	}, nil
}

func checkVirtualHardDiskExistsByPath(path string) (exists bool) {
	// 获取文件所在目录
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CreateVirtualHardDisk(path string, name string, sizeGiB int) (vhd *VirtualHardDisk, err error) {
	if checkVirtualHardDiskExistsByPath(path) {
		return nil, errors.New("VirtualHardDisk exists")
	}
	mgmt, err := wmictl.NewLocalImageManagementService()
	if err != nil {
		return
	}
	wHost := mgmt.GetWmiHost()
	vhdSettingData, err := disk.GetDefaultVirtualHardDiskSettingData(wHost)
	if err != nil {
		return
	}
	// Set the path of the disk
	if err = vhdSettingData.SetPropertyPath(path); err != nil {
		return
	}
	// Fixed 固定大小, Dynamic 动态硬盘, Differencing 差分硬盘
	// current default type is dynamic
	if err = vhdSettingData.SetPropertyType(disk.VirtualHardDiskType_SPARSE); err != nil {
		return nil, err
	}
	// VHD, VHDX, VHDSet
	// current default format is VHDX
	if err = vhdSettingData.SetPropertyFormat(uint16(disk.VirtualHardDiskFormat_2)); err != nil {
		return nil, err
	}
	if err = vhdSettingData.SetPropertyMaxInternalSize(uint64(sizeGiB) * 1024 * 1024 * 1024); err != nil {
		return
	}
	err = mgmt.CreateDisk(vhdSettingData)
	vhd = &VirtualHardDisk{
		Name:        name,
		Path:        path,
		UsedSizeGB:  0,
		TotalSizeGB: uint64(sizeGiB),
	}
	return
}
