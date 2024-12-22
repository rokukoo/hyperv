package wmictl

import "github.com/microsoft/wmi/pkg/constant"

type WmiNamespace = string

const (
	VirtualizationV2 WmiNamespace = string(constant.Virtualization)
)
