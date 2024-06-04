package test

import (
	"testing"

	datalink "github.com/Mayvid0/netSimGo/internal/dataLinkLayer"
	"github.com/Mayvid0/netSimGo/internal/network"
	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/utils"
)

func TestNetworkLayerFunctionalities(t *testing.T) {
	// Initialize devices for two different local networks
	devices1 := []*physical.Device{
		{Name: "Device1", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()},
		{Name: "Device2", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()},
	}

	devices2 := []*physical.Device{
		{Name: "Device3", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()},
		{Name: "Device4", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()},
	}

	// Create star topologies with switches for each network
	starTopology1 := &datalink.StarTopologyWithSwitch{}
	switchDevice1 := &physical.Switch{
		Hub: physical.Hub{
			Device: physical.Device{
				Name:       "Switch 1",
				MACAddress: utils.GenerateRandomMAC(),
				LinkStatus: true,
			},
			NumberPorts: 3,
		},
		SwitchingTable: make(map[string]int),
	}
	starTopology1.Switch = *switchDevice1

	starTopology2 := &datalink.StarTopologyWithSwitch{}
	switchDevice2 := &physical.Switch{
		Hub: physical.Hub{
			Device: physical.Device{
				Name:       "Switch 2",
				MACAddress: utils.GenerateRandomMAC(),
				LinkStatus: true,
			},
			NumberPorts: 3,
		},
		SwitchingTable: make(map[string]int),
	}
	starTopology2.Switch = *switchDevice2

	// Create router interfaces
	routerInterface1 := &physical.Device{Name: "eth1", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()}
	routerInterface2 := &physical.Device{Name: "eth2", LinkStatus: true, MACAddress: utils.GenerateRandomMAC()}

	// Add and connect devices to the star topologies and connect switches to router interfaces
	starTopology1.AddEndDevice(devices1[0])
	starTopology1.AddEndDevice(devices1[1])
	starTopology1.AddEndDevice(routerInterface1)
	starTopology1.ConnectEndDevice(switchDevice1, routerInterface1)
	starTopology1.ConnectEndDevice(switchDevice1, devices1[0])
	starTopology1.ConnectEndDevice(switchDevice1, devices1[1])

	starTopology2.AddEndDevice(routerInterface2)
	starTopology2.AddEndDevice(devices2[0])
	starTopology2.AddEndDevice(devices2[1])
	starTopology2.ConnectEndDevice(switchDevice2, routerInterface2)
	starTopology2.ConnectEndDevice(switchDevice2, devices2[0])
	starTopology2.ConnectEndDevice(switchDevice2, devices2[1])

	for _, device := range starTopology1.EndDevices {
		t.Logf("Assigned Mac address %s to  %s", device.MACAddress, device.Name)
	}
	for _, device := range starTopology2.EndDevices {
		t.Logf("Assigned Mac address %s to  %s", device.MACAddress, device.Name)
	}

	utils.AssignIPAddress(starTopology1.EndDevices, "192.168.1.0/24")
	utils.AssignIPAddress(starTopology2.EndDevices, "192.168.2.0/24")

	for _, device := range starTopology1.EndDevices {
		t.Logf("Assigned IP address %s to  %s", device.IPAddress, device.Name)
	}
	for _, device := range starTopology2.EndDevices {
		t.Logf("Assigned IP address %s to  %s", device.IPAddress, device.Name)
	}

	// Initialize the router's ARP table for testing
	router := network.Router{Interfaces: []physical.Device{*routerInterface1, *routerInterface2}}
	router.ARPTable = make(map[string]network.ARPEntry)

	t.Logf("Network topology created successfully")
}
