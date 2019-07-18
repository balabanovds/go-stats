package globs

// Element contains details about each network element
type Element struct {
	IP          string
	ElementType int
	Policy      Policy
	Server      string
	Err         error
}

// Policy struct contains details about each mediation entry
type Policy struct {
	ID          int
	SnmpVersion string
	Community   string
	UserName    string
}
