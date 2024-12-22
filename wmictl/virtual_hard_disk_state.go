package wmictl

import (
	"encoding/xml"
	"strconv"
)

// 定义结构体映射XML结构
type virtualHardDiskStateInstance struct {
	XMLName   xml.Name   `xml:"INSTANCE"`
	ClassName string     `xml:"CLASSNAME,attr"`
	Property  []property `xml:"PROPERTY"`
}

type property struct {
	XMLName    xml.Name `xml:"PROPERTY"`
	Name       string   `xml:"NAME,attr"`
	Type       string   `xml:"TYPE,attr"`
	Value      string   `xml:"VALUE"`
	Propagated string   `xml:"PROPAGATED,attr,omitempty"`
}

type VirtualHardDiskState struct {
	Alignment               uint32
	FileSize                uint64
	FragmentationPercentage uint32
	InUse                   bool
	MinInternalSize         uint64
	PhysicalSectorSize      uint32
	Timestamp               string
}

func (ins *virtualHardDiskStateInstance) virtualHardDiskState() (*VirtualHardDiskState, error) {
	vhdState := &VirtualHardDiskState{}
	for _, prop := range ins.Property {
		value := prop.Value
		switch prop.Name {
		case "Alignment":
			// 字符串值 转 uint32
			atoi, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			vhdState.Alignment = uint32(atoi)
		case "FileSize":
			// 字符串值 转 uint64
			atoi, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			vhdState.FileSize = uint64(atoi)
		case "FragmentationPercentage":
			// 字符串值 转 uint32
			atoi, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			vhdState.FragmentationPercentage = uint32(atoi)
		case "InUse":
			// 字符串值 转 bool
			atoi, err := strconv.ParseBool(value)
			if err != nil {
				return nil, err
			}
			vhdState.InUse = atoi
		case "MinInternalSize":
			if value == "" {
				continue
			}
			// 字符串值 转 uint64
			atoi, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			vhdState.MinInternalSize = uint64(atoi)
		case "PhysicalSectorSize":
			// 字符串值 转 uint32
			atoi, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			vhdState.PhysicalSectorSize = uint32(atoi)
		case "Timestamp":
			// 字符串值 转 string
			vhdState.Timestamp = value
		}
	}
	return vhdState, nil
}
