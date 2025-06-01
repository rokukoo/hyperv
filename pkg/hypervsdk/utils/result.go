package hypervsdk

import (
	"fmt"
	"github.com/rokukoo/hyperv/pkg/wmiext"
	"strings"
)

func WaitResult(res int32, service *wmiext.Service, job *wmiext.Instance, errorMsg string, translate func(int) error) error {
	var err error

	switch res {
	case 0:
		return nil
	case 4096:
		err = wmiext.WaitJob(service, job)
		//defer job.Close()
	default:
		if translate != nil {
			return translate(int(res))
		}

		return fmt.Errorf("%s (result code %d)", errorMsg, res)
	}

	if err != nil {
		desc, _ := job.GetAsString("ErrorDescription")
		desc = strings.Replace(desc, "\n", " ", -1)
		return fmt.Errorf("%s: %w (%s)", errorMsg, err, desc)
	}

	return err
}
