package hypervctl

import (
	"github.com/pkg/errors"
)

var (
	VirtualMachineAlreadyExists = errors.New("virtual machine already exists")
)
