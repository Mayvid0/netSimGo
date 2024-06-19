package physical

import "fmt"

type ARPEntry struct {
	MACAddress string
}

type Device struct {
	Name       string
	MACAddress string
	IPAddress  string
	LinkStatus bool
	PortNumber int
	HasToken   bool
	ARPTable   map[string]ARPEntry
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

func (d *Device) CreateArpRequest(destination *Device, s *Switch) {
	fmt.Printf("Sending ARP request from MAC %s for IP %s\n", d.MACAddress, destination.IPAddress)

	// Find the source port based on the source MAC address
	// sourcePort := s.FindSourcePort(d.MACAddress)

	// // Update the MAC address table with the source MAC address and port number
	sourcePort := d.PortNumber
	s.SwitchingTable[d.MACAddress] = sourcePort
	// sourcePort := d.PortNumber

	// Broadcast the ARP request to all end devices except the source port
	for _, endDevice := range s.EndDevices {
		fmt.Printf("Broadcasting ARP request to %s\n", endDevice.Name)

		if endDevice.PortNumber != sourcePort {
			fmt.Printf("Broadcasting ARP request to %s\n", endDevice.Name)
			endDevice.HandleArpRequest(d, destination, s)
		}
	}
}

func (d *Device) HandleArpRequest(source *Device, destination *Device, s *Switch) {
	if d.IPAddress == destination.IPAddress {
		d.addToArp(source)
		fmt.Printf("Sending ARP reply to MAC %s for IP %s\n", source.MACAddress, destination.IPAddress)
		d.SendReplytoSwitch(source, s)
	} else {
		fmt.Printf("Ignoring ARP request for IP %s\n", destination.IPAddress)
	}
}

func (d *Device) SendReplytoSwitch(desti *Device, s *Switch) {
	s.SwitchingTable[d.MACAddress] = d.PortNumber
	_, ok := s.SwitchingTable[d.MACAddress]

	if ok {
		s.SendDataFromSwitch(d, desti, d.MACAddress)
	} else {

		s.SwitchingTable[d.MACAddress] = d.PortNumber
		s.SendDataFromSwitch(d, desti, d.MACAddress)
	}
}

func (s *Switch) SendDataFromSwitch(source *Device, desti *Device, message string) {
	fmt.Printf("Sending result of arp query back to %s\n", desti.MACAddress)
	desti.addToArp(source)
}

func (d *Device) addToArp(addDevice *Device) {
	fmt.Printf("Adding device %s to %s arp table\n", addDevice.Name, d.Name)
	d.ARPTable[addDevice.IPAddress] = (ARPEntry{addDevice.MACAddress})
}
