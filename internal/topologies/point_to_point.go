package topologies

import (
	"fmt"

	"github.com/Mayvid0/netSimGo/internal/physical"
)

// basic struct for point-to-point topology
type PointToPoint struct {
	Topology
	Devices []*physical.Device
}

// defining the methods for the point-to-point topology

func (p *PointToPoint) AddEndDevice(device *physical.Device) {
	p.Devices = append(p.Devices, device)
}

func (p *PointToPoint) ConnectEndDevice(device1 *physical.Device, device2 *physical.Device) error {

	if !device1.LinkStatus || !device2.LinkStatus {
		return fmt.Errorf("cannot connect devices %s, %s and %s, %s: one or both devices are not linked", device1.Name, device1.MACAddress, device2.Name, device2.MACAddress)
	}

	connection := physical.Connection{Device1: device1, Device2: device2}
	fmt.Printf("Connection established between %s and %s\n", connection.Device1.Name, connection.Device2.Name)

	return nil
}

func (p *PointToPoint) ReceiveData(source *physical.Device, receiver *physical.Device, message string) {
	fmt.Printf("Received message from %s to %s: %s\n", source.Name, receiver.Name, message)
}

func (p *PointToPoint) SendData(source *physical.Device, receiver *physical.Device, message string) {
	fmt.Printf("Sending message from %s to %s: %s\n", source.Name, receiver.Name, message)
	p.ReceiveData(source, receiver, message)
}
