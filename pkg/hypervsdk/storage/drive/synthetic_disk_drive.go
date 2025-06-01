package drive

import "github.com/rokukoo/hyperv/pkg/wmiext"

type SyntheticDiskDrive struct {
	*VirtualDrive
}

// NewSyntheticDiskDrive creates a new SyntheticDiskDrive instance
func NewSyntheticDiskDrive(instance *wmiext.Instance) (*SyntheticDiskDrive, error) {
	vdriver, err := NewVirtualDrive(instance)
	if err != nil {
		return nil, err
	}
	return &SyntheticDiskDrive{vdriver}, nil
}
