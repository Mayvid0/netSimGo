package physical

// there should be a struct for device (doesnt matter what layer... a basic struct
// all other layer devices can inherit from it)

type Device struct {
	Name       string
	MACAddress string
	LinkStatus bool
	PortNumber int
	HasToken   bool
}

type Hub struct {
	Device
	NumberPorts int
	EndDevices  []*Device
}

type Connection struct {
	Device1 *Device
	Device2 *Device
}

type Bridge struct {
	Device
	NumberPorts     int
	EndDevices      []*Hub
	ForwardingTable map[string]int
}

type Switch struct {
	Hub
	SwitchingTable map[string]int
}
