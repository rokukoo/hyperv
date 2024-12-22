package hypervctl

import "github.com/pkg/errors"

var (
	ErrVmAlreadyExists    = errors.New("failed to create virtual machine: vm already exists")
	ErrVmAlreadyRunning   = errors.New("vm is already running")
	ErrVmAlreadySuspended = errors.New("vm is already suspended")
	ErrVmAlreadyStopped   = errors.New("vm is already stopped")
	ErrVmAlreadySaved     = errors.New("vm is already saved")
)
