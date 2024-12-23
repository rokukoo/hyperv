package hypervctl

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

var (
	ErrVmAlreadyExists    = errors.New("failed to create virtual machine: vm already exists")
	ErrVmAlreadyRunning   = errors.New("vm is already running")
	ErrVmAlreadySuspended = errors.New("vm is already suspended")
	ErrVmAlreadyStopped   = errors.New("vm is already stopped")
	ErrVmAlreadySaved     = errors.New("vm is already saved")
)

var (
	ErrHyperVNamespaceMissing = errors.New("HyperV namespace not found, is HyperV enabled?")
)

func translateCommonHyperVWmiError(wmiError error) error {
	if werr, ok := wmiError.(*wmiext.WmiError); ok {
		switch werr.Code() {
		case wmiext.WBEM_E_INVALID_NAMESPACE:
			return ErrHyperVNamespaceMissing
		}
	}

	return wmiError
}
