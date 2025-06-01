package hyperv

import (
	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("not found")

	ErrorVirtualMachineAlreadyExists = errors.New("virtual machine already exists")
)
