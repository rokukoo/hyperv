package networking

import (
	"fmt"
	"github.com/rokukoo/hypervctl/pkg/wmiext"
)

const (
	Msvm_ExternalEthernetPort = "Msvm_ExternalEthernetPort"
)

type ExternalEthernetPort struct {
	S__PATH  string `json:"-"`
	S__CLASS string `json:"-"`

	PermanentAddress string

	*wmiext.Instance
}

func (eep *ExternalEthernetPort) Path() string {
	return eep.S__PATH
}

func GetExternalEthernetPort(con *wmiext.Service, ethernetName string) (*ExternalEthernetPort, error) {
	extPort := &ExternalEthernetPort{}
	wquery := fmt.Sprintf("SELECT * FROM Msvm_ExternalEthernetPort WHERE ElementName = '%s'", ethernetName)
	return extPort, con.FindFirstObject(wquery, extPort)
}
