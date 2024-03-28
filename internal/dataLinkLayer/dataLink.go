package datalink

import (
	"fmt"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

// we must create a struct like star topology , but with switch

type StarTopologyWithSwitch struct {

	// import the topology functions and the switch struct created in physical layer
	physical.Switch
	topologies.Topology
}

func (s *StarTopologyWithSwitch) AddEndDevice(device *physical.Device) {
	s.EndDevices = append(s.EndDevices, device)
}

func (s *StarTopologyWithSwitch) ConnectEndDevice(switchDevice *physical.Switch, device *physical.Device) error {
	if !switchDevice.LinkStatus || !device.LinkStatus {
		return fmt.Errorf("cannot connect devices %s and %s: one or both devices are not linked", switchDevice.Name, device.Name)
	}

	fmt.Printf("Connection established between switch %s and end device %s\n", switchDevice.Name, device.Name)
	return nil
}

func (s *StarTopologyWithSwitch) SwitchingTable(switchDevice *physical.Switch) {
	// this would create a table of the mapping of each port to its corresponding mac address
	//performs address learning
	for _, device := range s.EndDevices {

		// go through all endDevices and then map macddress to their portNumber
		switchDevice.SwitchingTable[device.MACAddress] = device.PortNumber
	}

}

func (s *StarTopologyWithSwitch) SendDataFromSwitch(port int, message string) {
	// Find the end device with the specified port number
	var targetDevice *physical.Device
	for _, device := range s.EndDevices {
		if device.PortNumber == port {
			targetDevice = device
			break
		}
	}

	if targetDevice != nil {
		fmt.Printf("Sending message from Switch %s to %s: %s\n", s.Name, targetDevice.Name, message)
		s.ReceiveDataFromHub(targetDevice, message)
		// Send the message to the target device here
	} else {
		fmt.Println("Error: End device not found for port number", port)
	}
}

func (s *StarTopologyWithSwitch) ReceiveDataFromHub(receiver *physical.Device, message string) {
	fmt.Printf(" Device %s Received message:  %s \n", receiver.Name, message)
}

func (s *StarTopologyWithSwitch) SendDataToSwitch(switchDevice *physical.Switch, source *physical.Device, receiver *physical.Device, message string) {
	// we send receiver's mac address along with the message
	// if we dont find the receiver mac address in switching table, we perform address learning
	port, ok := switchDevice.SwitchingTable[receiver.MACAddress]
	if ok {
		fmt.Printf("Sending message from %s to Switch %s: %s\n", source.Name, switchDevice.Name, message)
		s.SendDataFromSwitch(port, message)
	} else {
		//perform address learning
		s.SwitchingTable(switchDevice)
		s.SendDataToSwitch(switchDevice, source, receiver, message)
	}
}
