package errors

import (
	"github.com/pkg/errors"
	"github.com/rokukoo/hyperv/pkg/wmiext"
)

var (
	ErrHyperVNamespaceMissing = errors.New("HyperV namespace not found, is HyperV enabled?")
)

func TranslateCommonHyperVWmiError(wmiError error) error {
	if werr, ok := wmiError.(*wmiext.WmiError); ok {
		switch werr.Code() {
		case wmiext.WBEM_E_INVALID_NAMESPACE:
			return ErrHyperVNamespaceMissing
		}
	}

	return wmiError
}
