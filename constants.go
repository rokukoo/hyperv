package hypervctl

// type VMStatus = virtualsystem.VirtualMachineState
type VMStatus = int

const (
	VMStatusUnknown    VMStatus = 0
	VMStatusCreating            = 10
	VMStatusStarting            = 20
	VMStatusRunning             = 30
	VMStatusStopping            = 40
	VMStatusStopped             = 50
	VMStatusSuspending          = 60
	VMStatusSuspended           = 70
	VMStatusSaving              = 80
	VMStatusSaved               = 90
	VMStatusRebooting           = 100
)
