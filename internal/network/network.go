package network

import "github.com/Mayvid0/netSimGo/internal/physical"

type RoutingEntry struct {
	NetworkAddressMask string
	NextHop            string
	Interface          string
}

type ARPEntry struct {
	IPAddress  string
	MACAddress string
}

type Router struct {
	ID              string
	Interfaces      []physical.Device
	RoutingTable    []RoutingEntry
	ARPTable        map[string]ARPEntry
	RoutingProtocol string
	RIPData         *RIPData
	OSPFData        *OSPFData
}

// RIPData and OSPFData structures would be defined according to the needs of the specific protocols.
type RIPData struct {
	// Specific fields for RIP protocol
}

type OSPFData struct {
	// Specific fields for OSPF protocol
}
