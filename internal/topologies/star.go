package topologies

import (
	"fmt"

	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
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
	Timestamp      time.Time
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

func (s *Star) SendDataToHub(source *physical.Device, hub *physical.Hub, packets []Packet) {
	fmt.Printf("Sending packets from %s to hub %s\n", source.Name, hub.Name)
	for _, packet := range packets {
		s.SendPacketFromHub(source, hub, packet)
	}
}

func (s *Star) SendPacketFromHub(source *physical.Device, hub *physical.Hub, packet Packet) {
	for _, device := range s.EndDevices {
		if device.MACAddress == source.MACAddress {
			continue
		}
		fmt.Printf("Broadcast packet from hub %s to %s\n", hub.Name, device.Name)
		s.ReceivePacketFromHub(source, device, packet)
	}
}

func (s *Star) ReceivePacketFromHub(source *physical.Device, receiver *physical.Device, packet Packet) {
	fmt.Printf("Device %s received packet from %s via hub\n", receiver.Name, source.Name)

	// Calculate checksum
	receivedChecksum := calculateChecksum(packet.Data)
	if receivedChecksum != packet.Checksum {
		fmt.Println("Checksum mismatch. Packet may be corrupted.")
		return
	}

	// Convert the packet data to a string
	receivedString := string(packet.Data)

	// Display the received string
	fmt.Println("Received data: ", receivedString)
}

// Calculate checksum for a byte array
func calculateChecksum(data []byte) uint16 {
	var sum uint32

	// Sum all 16-bit words in the data
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}

	// If the data length is odd, add the last byte as a 16-bit word
	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}

	// Fold the 32-bit sum to 16 bits
	for sum>>16 > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	// Take the one's complement of the sum
	checksum := uint16(^sum)

	return checksum
}
