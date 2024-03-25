package topologies

import (
	"fmt"

	"github.com/Mayvid0/netSimGo/internal/physical"
)

type Star struct {
	physical.Hub
	Topology
}

//define all the methods of the topology

func (s *Star) AddEndDevice(device *physical.Device) {
	s.EndDevices = append(s.EndDevices, device)
}

func (s *Star) ConnectEndDevice(hub *physical.Hub, device *physical.Device) error {
	if !hub.LinkStatus || !device.LinkStatus {
		return fmt.Errorf("cannot connect devices %s and %s: one or both devices are not linked", hub.Name, device.Name)
	}
	// Here you might perform additional connection logic specific to the Star topology
	fmt.Printf("Connection established between hub %s and end device %s\n", hub.Name, device.Name)
	return nil
}

func (s *Star) SendDataToHub(source *physical.Device, hub *physical.Hub, message string) {
	fmt.Printf("Sending message from %s to hub %s: %s\n", source.Name, hub.Name, message)
	s.SendDataFromHub(source, hub, message)
}

func (s *Star) SendDataFromHub(source *physical.Device, hub *physical.Hub, message string) {
	for _, device := range s.EndDevices {
		if device.MACAddress == source.MACAddress {
			continue
		}
		fmt.Printf("Broadcast message from hub %s to %s: %s\n", hub.Name, device.Name, message)
		s.ReceiveDataFromHub(source, device, message)
	}

}

func (s *Star) ReceiveDataFromHub(source *physical.Device, receiver *physical.Device, message string) {
	fmt.Printf(" Device %s Received message %s broadcasted through hub : %s\n", receiver.Name, source.Name, message)
}
