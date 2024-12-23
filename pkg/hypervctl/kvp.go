package hypervctl

type CimKvpItemProperty struct {
	Name  string `xml:"NAME,attr"`
	Value string `xml:"VALUE"`
}

type KvpError struct {
	ErrorCode int
	message   string
}
