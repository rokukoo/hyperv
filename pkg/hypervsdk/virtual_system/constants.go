package virtual_system

type ComputerSystemState int32

// https://docs.microsoft.com/en-us/previous-versions/windows/desktop/virtual/msvm-computersystem?redirectedfrom=MSDN
const (
	Unknown            ComputerSystemState = 0
	Other              ComputerSystemState = 1
	Running            ComputerSystemState = 2
	Off                ComputerSystemState = 3
	Stopping           ComputerSystemState = 4
	Saved              ComputerSystemState = 6
	Paused             ComputerSystemState = 9
	Starting           ComputerSystemState = 10
	Reset              ComputerSystemState = 11
	Saving             ComputerSystemState = 32773
	Pausing            ComputerSystemState = 32776
	Resuming           ComputerSystemState = 32777
	FastSaved          ComputerSystemState = 32779
	FastSaving         ComputerSystemState = 32780
	ForceShutdown      ComputerSystemState = 32781
	ForceReboot        ComputerSystemState = 32782
	Hibernated         ComputerSystemState = 32783
	ComponentServicing ComputerSystemState = 32784
	RunningCritical    ComputerSystemState = 32785
	OffCritical        ComputerSystemState = 32786
	StoppingCritial    ComputerSystemState = 32787
	SavedCritical      ComputerSystemState = 32788
	PausedCritical     ComputerSystemState = 32789
	StartingCritical   ComputerSystemState = 32790
	ResetCritical      ComputerSystemState = 32791
	SavingCritical     ComputerSystemState = 32792
	PausingCritical    ComputerSystemState = 32793
	ResumingCritical   ComputerSystemState = 32794
	FastSaveCritical   ComputerSystemState = 32795
	FastSavingCritical ComputerSystemState = 32796
)

const (
	StateChangeTimeoutSeconds = 300
)
