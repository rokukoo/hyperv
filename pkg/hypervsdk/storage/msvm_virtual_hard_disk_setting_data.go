package storage

import (
	"encoding/xml"
	utils "github.com/rokukoo/hyperv/pkg/hypervsdk/utils"
	"github.com/rokukoo/hyperv/pkg/wmiext"
	"strconv"
	"time"
)

const (
	Msvm_VirtualHardDiskSettingData = "Msvm_VirtualHardDiskSettingData"
)

type VirtualHardDiskType uint16

const (
	VirtualHardDiskType_NONE   = 0
	VirtualHardDiskType_LEGACY = 1
	VirtualHardDiskType_FLAT   = 2
	VirtualHardDiskType_SPARSE = 3
)

// VHD, VHDX, VHDSet
type VirtualHardDiskFormat uint16

const (
	VirtualHardDiskFormat_NONE = 0
	VirtualHardDiskFormat_ISO  = 1
	VirtualHardDiskFormat_1    = 2
	VirtualHardDiskFormat_2    = 3
)

type VirtualHardDiskSettingData struct {
	InstanceID                 string    `json:"instance_id"`
	Caption                    string    `json:"caption"`
	Description                string    `json:"description"`
	ElementName                string    `json:"element_name"`
	Type                       uint16    `json:"type"`
	Format                     uint16    `json:"format"`
	Path                       string    `json:"path"`
	ParentPath                 string    `json:"parent_path"`
	ParentTimestamp            time.Time `json:"parent_timestamp"`
	ParentIdentifier           string    `json:"parent_identifier"`
	MaxInternalSize            uint64    `json:"max_internal_size"`
	BlockSize                  uint32    `json:"block_size"`
	LogicalSectorSize          uint32    `json:"logical_sector_size"`
	PhysicalSectorSize         uint32    `json:"physical_sector_size"`
	VirtualDiskId              string    `json:"virtual_disk_id"`
	DataAlignment              uint64    `json:"data_alignment"`
	PmemAddressAbstractionType uint16    `json:"pmem_address_abstraction_type"`
	IsPmemCompatible           bool      `json:"is_pmem_compatible"`

	Size        uint64 `json:"size"`
	LSectorSize uint32 `json:"l_sector_size"`
	PSectorSize uint32 `json:"p_sector_size"`

	//service *virtualsystem.VirtualSystemManagementService `json:"-"`
	*wmiext.Instance
}

func (ims *ImageManagementService) GetDefaultVirtualHardDiskSettingData() (*VirtualHardDiskSettingData, error) {
	var (
		vhsd = VirtualHardDiskSettingData{}
		err  error
	)
	if vhsd.Instance, err = ims.Session.CreateInstance(Msvm_VirtualHardDiskSettingData, &vhsd); err != nil {
		return nil, err
	}
	return &vhsd, nil
}

func (ims *ImageManagementService) NewVirtualHardDiskSettingData(
	path string,
	logicalSectorSize, physicalSectorSize, blockSize uint32,
	diskSize uint64,
	dynamic bool,
	diskFileFormat VirtualHardDiskFormat,
) (
	vhdsetting *VirtualHardDiskSettingData,
	err error,
) {
	vhdsetting, err = ims.GetDefaultVirtualHardDiskSettingData()
	if err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("Path", path); err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("Format", uint16(diskFileFormat)); err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("BlockSize", blockSize); err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("LogicalSectorSize", logicalSectorSize); err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("PhysicalSectorSize", physicalSectorSize); err != nil {
		return nil, err
	}

	if err = vhdsetting.Put("MaxInternalSize", diskSize); err != nil {
		return nil, err
	}

	// Fixed 固定大小, Dynamic 动态硬盘, Differencing 差分硬盘
	// current default type is dynamic
	if dynamic {
		if err = vhdsetting.Put("Type", uint16(VirtualHardDiskType_SPARSE)); err != nil {
			return nil, err
		}
	} else {
		// Fixed
		if err = vhdsetting.Put("Type", uint16(VirtualHardDiskType_FLAT)); err != nil {
			return nil, err
		}
	}

	return vhdsetting, nil
}

type INSTANCE struct {
	XMLName   xml.Name `xml:"INSTANCE"`
	Text      string   `xml:",chardata"`
	CLASSNAME string   `xml:"CLASSNAME,attr"`
	PROPERTY  []struct {
		Text       string `xml:",chardata"`
		NAME       string `xml:"NAME,attr"`
		TYPE       string `xml:"TYPE,attr"`
		PROPAGATED string `xml:"PROPAGATED,attr"`
		VALUE      string `xml:"VALUE"`
	} `xml:"PROPERTY"`
}

func getVirtualHardDiskSettingDataFromXml(
	xmlInstance string,
) (
	*VirtualHardDiskSettingData,
	error,
) {
	var virtualHardDiskSettingData = &VirtualHardDiskSettingData{}
	var err error
	var t uint64
	var diskData INSTANCE

	if err = xml.Unmarshal([]byte(xmlInstance), &diskData); err != nil {
		return nil, err
	}

	for _, property := range diskData.PROPERTY {
		switch property.NAME {
		case "MaxInternalSize":
			virtualHardDiskSettingData.Size, err = strconv.ParseUint(property.VALUE, 10, 64)
			if err != nil {
				return nil, err
			}
		case "BlockSize":
			t, err = strconv.ParseUint(property.VALUE, 10, 32)
			if err != nil {
				return nil, err
			}
			virtualHardDiskSettingData.BlockSize = uint32(t)
		case "LogicalSectorSize":
			t, err = strconv.ParseUint(property.VALUE, 10, 32)
			if err != nil {
				return nil, err
			}
			virtualHardDiskSettingData.LogicalSectorSize = uint32(t)
		case "PhysicalSectorSize":
			t, err = strconv.ParseUint(property.VALUE, 10, 32)
			if err != nil {
				return nil, err
			}
			virtualHardDiskSettingData.PhysicalSectorSize = uint32(t)
		case "Format":
			t, err = strconv.ParseUint(property.VALUE, 10, 16)
			if err != nil {
				return nil, err
			}
			virtualHardDiskSettingData.Format = uint16(t)
		case "ParentPath":
			virtualHardDiskSettingData.ParentPath = property.VALUE
		}
	}

	return virtualHardDiskSettingData, nil
}

func (ims *ImageManagementService) GetVirtualHardDiskSettingData(path string) (*VirtualHardDiskSettingData, error) {
	var (
		err         error
		job         *wmiext.Instance
		settingData string
		returnValue int32
	)

	if err = ims.Method("GetVirtualHardDiskSettingData").
		In("Path", path).
		Execute().
		Out("SettingData", &settingData).
		Out("Job", &job).
		Out("ReturnValue", &returnValue).
		End(); err != nil {
		return nil, err
	}

	if err = utils.WaitResult(returnValue, ims.Session, job, "Failed to get setting data for disk", nil); err != nil {
		return nil, err
	}

	return getVirtualHardDiskSettingDataFromXml(settingData)
}
