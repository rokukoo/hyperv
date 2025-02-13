package hypervctl

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

var (
	VirtualMachineAlreadyExists = errors.New("virtual machine already exists")
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
