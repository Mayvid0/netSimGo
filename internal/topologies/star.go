package topologies

import (
	"fmt"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/utils"
)

type Star struct {
	physical.Hub
	Topology
}

type Packet struct {
	SequenceNumber int
	Acknowledgment bool
	IsEnd          bool
	Data           []byte
	Checksum       uint16
	Retransmit     bool // Flag to indicate if the packet needs retransmission
}

// Define all the methods of the topology

func (s *Star) AddEndDevice(device *physical.Device) {
	s.EndDevices = append(s.EndDevices, device)
}

func (s *Star) ConnectEndDevice(hub *physical.Hub, device *physical.Device) error {
	if !hub.LinkStatus || !device.LinkStatus {
		return fmt.Errorf("cannot connect devices %s and %s: one or both devices are not linked", hub.Name, device.Name)
	}

	fmt.Printf("Connection established between hub %s and end device %s\n", hub.Name, device.Name)
	return nil
}

func (s *Star) SendDataToHub(source *physical.Device, hub *physical.Hub, packets []Packet) int {
	fmt.Printf("Sending frames from %s to hub %s\n", source.Name, hub.Name)
	for _, packet := range packets {
		s.SendPacketFromHub(source, hub, packet)
	}
	return len(packets)
}

func (s *Star) SendPacketFromHub(source *physical.Device, hub *physical.Hub, packet Packet) {
	for _, device := range s.EndDevices {
		if device.MACAddress == source.MACAddress {
			continue
		}
		fmt.Printf("Broadcast frame from hub %s to %s\n", hub.Name, device.Name)
		s.ReceivePacketFromHub(source, device, packet)
	}
}

func (s *Star) ReceivePacketFromHub(source *physical.Device, receiver *physical.Device, packet Packet) {
	fmt.Printf("Device %s received frame from %s via hub\n", receiver.Name, source.Name)

	// Calculate checksum
	receivedChecksum := utils.CalculateChecksum(packet.Data)
	if receivedChecksum != packet.Checksum {
		fmt.Println("Checksum mismatch. frame may be corrupted.")
		return
	}

	// Convert the packet data to a string
	receivedString := string(packet.Data)

	// Display the received string
	fmt.Println("Received data: ", receivedString)
}

func (s *Star) SendDataToBridge(source *physical.Hub, destination *physical.Hub, bridge *physical.Bridge, packets []Packet) int {
	fmt.Printf("Sending frames from %s to bridge %s\n", source.Name, bridge.Name)
	for _, packet := range packets {
		s.SendPacketToBridge(source, destination, bridge, packet)
	}
	return len(packets)
}

func (s *Star) SendPacketToBridge(source *physical.Hub, destination *physical.Hub, bridge *physical.Bridge, packet Packet) {
	fmt.Printf("Sending packet from %s to bridge %s\n", source.Name, bridge.Name)
	s.ForwardPacket(packet, source, destination) // Send packet to bridge with source hub and no destination hub
}

func (s *Star) ForwardPacket(packet Packet, sourceHub *physical.Hub, destinationHub *physical.Hub) {
	fmt.Printf("Forwarding packet from hub %s to hub %s via bridge \n", sourceHub.Name, destinationHub.Name)
	fmt.Println("Packet forwarded successfully.")
	s.SendHubData(sourceHub.EndDevices[0], destinationHub, packet)
}

func (s *Star) SendHubData(source *physical.Device, desti *physical.Hub, packet Packet) {
	for _, device := range desti.EndDevices {
		if device.MACAddress == source.MACAddress {
			continue
		}
		fmt.Printf("Broadcast frame from device %s to %s broadcasted via hub %s\n", source.Name, device.Name, desti.Name)
		s.ReceivePacketFromHub(source, device, packet)
	}
}
