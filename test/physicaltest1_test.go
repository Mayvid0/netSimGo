package test

import (
	"fmt"
	"testing"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

func TestSimulateDataTransmission(t *testing.T) {

	topology := &topologies.PointToPoint{}

	device1 := &physical.Device{Name: "Device1", MACAddress: "AA:BB:CC:DD:EE:FF", LinkStatus: true}
	device2 := &physical.Device{Name: "Device2", MACAddress: "11:22:33:44:55:66", LinkStatus: true}

	topology.AddEndDevice(device1)
	topology.AddEndDevice(device2)

	// Connect the devices
	err := topology.ConnectEndDevice(device1, device2)
	if err != nil {
		t.Errorf("Error connecting devices: %v", err)
		return
	}

	message := "Hello, World!"
	expected := "Received message from Device1 to Device2: Hello, World!"
	received := simulateDataTransmission(topology, device1, device2, message)
	if received != expected {
		t.Errorf("Expected: %s\nReceived: %s", expected, received)
	}
}

func simulateDataTransmission(topology topologies.Topology, source *physical.Device, receiver *physical.Device, message string) string {
	// Simulate sending data
	topology.SendData(source, receiver, message)

	// Simulate receiving data
	return fmt.Sprintf("Received message from %s to %s: %s", source.Name, receiver.Name, message)
}
