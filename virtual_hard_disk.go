package hyperv

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/resource"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/storage"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/storage/disk"
	"github.com/rokukoo/hyperv/pkg/hypervsdk/virtual_system"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

type VirtualHardDiskType = int

const (
	VirtualHardDiskTypeSystem = iota
	VirtualHardDiskTypeData
)

// VirtualHardDisk
// https://learn.microsoft.com/zh-cn/windows/win32/hyperv_v2/msvm-virtualharddisksettingdata
type VirtualHardDisk struct {
	Name        string              `json:"name"`
	Type        VirtualHardDiskType `json:"type"`
	TotalSizeGB float64             `json:"total_size"`
	UsedSizeGB  float64             `json:"used_size"`
	Path        string              `json:"path"`
	Attached    bool                `json:"attached"`
	*disk.VirtualHardDisk
}

func fileName(path string) string {
	// 根据路径读取文件名
	return filepath.Base(path)
}

func fileSizeGB(path string) (sizeGB float64, err error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	// 获取文件大小 (以 GB 为单位)
	return float64(fileInfo.Size()) / (1024 * 1024 * 1024), nil
}

func existsVirtualHardDiskByPath(path string) (exists bool) {
	// 获取文件所在目录
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func getVirtualHardDiskMaxSize(path string) (maxSizeGiB uint64, err error) {
	var ims *storage.ImageManagementService
	var vhdSettingData *storage.VirtualHardDiskSettingData
	if ims, err = storage.LocalImageManagementService(); err != nil {
		return
	}
	if vhdSettingData, err = ims.GetVirtualHardDiskSettingData(path); err != nil {
		return
	}
	return vhdSettingData.Size / (1024 * 1024 * 1024), nil
}

func (vhd *VirtualHardDisk) Detach() (err error) {
	if !vhd.Attached {
		return errors.New("vhd not attached")
	}
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return
	}
	if err = vmms.DetachVirtualHardDisk(vhd.VirtualHardDisk); err != nil {
		return
	}
	vhd.VirtualHardDisk = nil
	vhd.Attached = false
	return nil
}

func (vm *VirtualMachine) GetVirtualHardDisks() ([]*VirtualHardDisk, error) {
	var virtualHardDisks []*VirtualHardDisk
	settingData, err := vm.computerSystem.GetVirtualSystemSettingData()
	if err != nil {
		return nil, err
	}
	storageAllocationSettingDatas, err := settingData.GetStorageAllocationSettingData()
	if err != nil {
		return nil, err
	}
	for _, storageAllocationSettingData := range storageAllocationSettingDatas {
		hardDisk := disk.VirtualHardDisk{StorageAllocationSettingData: storageAllocationSettingData}
		path := hardDisk.HostResource[0]
		virtualHardDisk, err := GetVirtualHardDiskByPath(path)
		if err != nil {
			return nil, err
		}
		drive, err := hardDisk.GetDrive()
		if err != nil {
			return nil, err
		}
		controller, err := drive.GetParenObject()
		if err != nil {
			return nil, err
		}
		if controller.ResourceType == uint16(resource.ResourcePool_ResourceType_IDE_Controller) {
			virtualHardDisk.Type = VirtualHardDiskTypeSystem
		} else if controller.ResourceType == uint16(resource.ResourcePool_ResourceType_Parallel_SCSI_HBA) {
			virtualHardDisk.Type = VirtualHardDiskTypeData
		} else {
			return nil, errors.New("unknown controller type")
		}
		virtualHardDisks = append(virtualHardDisks, virtualHardDisk)
	}
	return virtualHardDisks, nil
}

func (vhd *VirtualHardDisk) Create() (err error) {
	var (
		mgmt *storage.ImageManagementService
	)
	if existsVirtualHardDiskByPath(vhd.Path) {
		return errors.New("VirtualHardDisk exists")
	}
	if mgmt, err = storage.LocalImageManagementService(); err != nil {
		return
	}
	vhdSettingData, err := mgmt.NewVirtualHardDiskSettingData(vhd.Path, 512, 512, 0, 1024*1024*1024*uint64(vhd.TotalSizeGB), true, storage.VirtualHardDiskFormat_2)
	if err != nil {
		return
	}
	if err = mgmt.CreateVirtualHardDisk(vhdSettingData); err != nil {
		return
	}
	vhd.Name = fileName(vhd.Path)
	if vhd.UsedSizeGB, err = fileSizeGB(vhd.Path); err != nil {
		return
	}
	//if vhd.VirtualHardDisk, err = disk.GetVirtualHardDiskByPath(mgmt.Session, vhd.Path); err != nil {
	//	return
	//}
	//vhd.Attached = true
	return nil
}

// AttachToByName 根据虚拟机名称挂载虚拟硬盘
// 
// 参数:
//   vmName: 虚拟机名称
// 返回:
//   bool: 是否挂载成功
//   error: 错误
func (vhd *VirtualHardDisk) AttachToByName(vmName string) (ok bool, err error) {
	virtualMachine, err := FirstVirtualMachineByName(vmName)
	if err != nil {
		return false, err
	}
	return vhd.AttachTo(virtualMachine)
}

func (vhd *VirtualHardDisk) AttachTo(virtualMachine *VirtualMachine) (ok bool, err error) {
	return vhd.AttachAsDataDisk(virtualMachine)
}

func (vhd *VirtualHardDisk) AttachAsDataDisk(virtualMachine *VirtualMachine) (ok bool, err error) {
	var (
		controllers []*resource.ResourceAllocationSettingData
	)
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	if controllers, err = virtualMachine.computerSystem.GetSCSIControllers(); errors.Is(err, wmiext.NotFound) || len(controllers) == 0 {
		if err = vmms.AddSCSIController(virtualMachine.computerSystem); err != nil {
			return false, err
		}
	} else if err != nil {
		return false, err
	}
	vhd.VirtualHardDisk, _, err = vmms.AttachVirtualHardDisk(virtualMachine.computerSystem, vhd.Path, virtual_system.VirtualHardDiskType_DATADISK_VIRTUALHARDDISK)
	if err != nil {
		return
	}
	vhd.Attached = true
	return true, nil
}

func (vhd *VirtualHardDisk) AttachAsSystemDisk(virtualMachine *VirtualMachine) (ok bool, err error) {
	vmms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return false, err
	}
	vhd.VirtualHardDisk, _, err = vmms.AttachVirtualHardDisk(virtualMachine.computerSystem, vhd.Path, virtual_system.VirtualHardDiskType_OS_VIRTUALHARDDISK)
	if err != nil {
		return
	}
	vhd.Attached = true
	return true, nil
}

// Resize 调整虚拟硬盘大小
// 
// 参数:
//   newSizeGiB: 新的虚拟硬盘大小 (GiB)
// 返回:
//   bool: 是否调整成功
//   error: 错误
func (vhd *VirtualHardDisk) Resize(newSizeGiB float64) (ok bool, err error) {
	ims, err := storage.LocalImageManagementService()
	if err != nil {
		return false, err
	}
	// TODO: check if there are snapshots, if there are, then will resize the snapshot
	//if snapshots, err := ims.GetSnapshotVirtualHardDisks(vhd.VirtualHardDisk); err != nil {
	//	return false, err
	//} else if len(snapshots) > 0 {
	//	if len(snapshots) == 1 {
	//
	//	}
	//	return false, errors.New("cannot resize a virtual hard disk that has many snapshots")
	//}

	if newSizeGiB <= vhd.UsedSizeGB {
		return false, errors.New("new size must be greater than used size")
	}
	err = ims.ResizeVirtualHardDisk(vhd.Path, uint64(newSizeGiB)*1024*1024*1024)
	if err != nil {
		return false, err
	}
	vhd.TotalSizeGB = newSizeGiB
	return true, nil
}

// GetVirtualHardDiskByPath 根据路径获取虚拟硬盘
// 
// 参数:
//   path: 虚拟硬盘路径
// 返回:
//   *VirtualHardDisk: 虚拟硬盘
//   error: 错误
func GetVirtualHardDiskByPath(path string) (*VirtualHardDisk, error) {
	if !existsVirtualHardDiskByPath(path) {
		return nil, errors.New("vhd not exists")
	}
	var virtualHardDisk *disk.VirtualHardDisk

	vsms, err := virtual_system.LocalVirtualSystemManagementService()
	if err != nil {
		return nil, err
	}

	// 获取文件大小 (以 GB 为单位)
	usedSizeGB, err := fileSizeGB(path)
	if err != nil {
		return nil, err
	}

	maxSizeGB, err := getVirtualHardDiskMaxSize(path)
	if err != nil {
		return nil, err
	}

	if virtualHardDisk, err = disk.GetVirtualHardDiskByPath(vsms.Session, path); err != nil {
		if !errors.Is(err, wmiext.NotFound) {
			return nil, err
		}
	}

	return &VirtualHardDisk{
		Name:            fileName(path),
		Path:            path,
		UsedSizeGB:      usedSizeGB,
		TotalSizeGB:     float64(maxSizeGB),
		Attached:        virtualHardDisk != nil,
		VirtualHardDisk: virtualHardDisk,
	}, nil
}

// CreateVirtualHardDisk 创建虚拟硬盘
// 
// 参数:
//   path: 虚拟硬盘路径
//   sizeGiB: 虚拟硬盘大小 (GiB)
// 返回:
//   *VirtualHardDisk: 虚拟硬盘
//   error: 错误
func CreateVirtualHardDisk(path string, sizeGiB float64) (vhd *VirtualHardDisk, err error) {
	vhd = &VirtualHardDisk{
		Path:        path,
		TotalSizeGB: sizeGiB,
	}
	if err = vhd.Create(); err != nil {
		return nil, err
	}
	return vhd, nil
}

// DeleteVirtualHardDiskByPath 根据路径删除虚拟硬盘
// 
// 参数:
//   path: 虚拟硬盘路径
// 返回:
//   bool: 是否删除成功
//   error: 错误
func DeleteVirtualHardDiskByPath(path string) (ok bool, err error) {
	if !existsVirtualHardDiskByPath(path) {
		return false, errors.New("vhd not exists")
	}
	if err = os.RemoveAll(path); err != nil {
		return
	}
	return true, nil
}
