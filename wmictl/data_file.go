package wmictl

import (
	"github.com/microsoft/wmi/pkg/base/host"
	"github.com/microsoft/wmi/pkg/base/query"
	"github.com/microsoft/wmi/pkg/constant"
	"github.com/microsoft/wmi/server2019/root/cimv2"
	"strings"
)

// SELECT FileSize, Name FROM CIM_DataFile WHERE Name = "D:\\Hyper-V\\Virtual Hard Disks\\新建虚拟硬盘.vhdx"

// FindDataFileByPath find data file by path
func FindDataFileByPath(name string) (*cimv2.CIM_DataFile, error) {
	whost := host.NewWmiLocalHost()
	wquery := query.NewWmiQuery("CIM_DataFile", "Name", strings.ReplaceAll(name, `\`, `\\`))
	creds := whost.GetCredential()
	return cimv2.NewCIM_DataFileEx6(whost.HostName, string(constant.CimV2), creds.UserName, creds.Password, creds.Domain, wquery)
}
