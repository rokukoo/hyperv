package storage

import (
	"github.com/rokukoo/hyperv/pkg/hypervsdk/storage/disk"
	utils "github.com/rokukoo/hyperv/pkg/hypervsdk/utils"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

const (
	Msvm_ImageManagementService = "Msvm_ImageManagementService"
)

type ImageManagementService struct {
	Session *wmiext.Service
	*wmiext.Instance
}

func LocalImageManagementService() (*ImageManagementService, error) {
	var (
		session *wmiext.Service
		svc     *wmiext.Instance
		err     error
	)
	// Get the WMI service
	if session, err = utils.NewLocalHyperVService(); err != nil {
		return nil, err
	}
	// Get the singleton instance
	if svc, err = session.GetSingletonInstance(Msvm_ImageManagementService); err != nil {
		return nil, err
	}
	return &ImageManagementService{session, svc}, nil
}

func (ims *ImageManagementService) ResizeVirtualHardDisk(path string, size uint64) (err error) {
	var (
		job         *wmiext.Instance
		returnValue int32
	)

	if err = ims.Method("ResizeVirtualHardDisk").
		In("Path", path).
		In("MaxInternalSize", size).
		Execute().
		Out("Job", &job).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return err
	}

	if err = utils.WaitResult(returnValue, ims.Session, job, "Failed to resize virtual hard disk", nil); err != nil {
		return err
	}

	return nil
}

func (ims *ImageManagementService) CreateVirtualHardDisk(settings *VirtualHardDiskSettingData) error {
	var (
		settingsObj string = settings.GetCimText()

		err         error
		job         *wmiext.Instance
		returnValue int32
	)

	if err = ims.Method("CreateVirtualHardDisk").
		In("VirtualDiskSettingData", settingsObj).
		Execute().
		Out("Job", &job).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return err
	}

	if err = utils.WaitResult(returnValue, ims.Session, job, "Failed to create virtual hard disk", nil); err != nil {
		return err
	}

	return nil
}

func (ims *ImageManagementService) GetSnapshotVirtualHardDisks(
	virtualHardDisk *disk.VirtualHardDisk,
) (
	snapshots []*disk.VirtualHardDisk,
	err error,
) {
	var virtualHardDiskSettingData *VirtualHardDiskSettingData
	virtualHardDisks := []*disk.VirtualHardDisk{}
	wquery := "SELECT * FROM Msvm_StorageAllocationSettingData"
	if err = ims.Session.FindObjects(wquery, virtualHardDisks); err != nil {
		return
	}
	for _, vhd := range virtualHardDisks {
		if virtualHardDiskSettingData, err = ims.GetVirtualHardDiskSettingData(vhd.GetPath()); err != nil {
			return nil, err
		}
		if virtualHardDiskSettingData.ParentPath != virtualHardDisk.GetPath() {
			continue
		}
		snapshots = append(snapshots, vhd)
	}
	return
}
