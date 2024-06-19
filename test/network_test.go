package test

import (
	"fmt"
	"testing"
	"time"

	datalink "github.com/Mayvid0/netSimGo/internal/dataLinkLayer"
	"github.com/Mayvid0/netSimGo/internal/network"
	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
	"github.com/Mayvid0/netSimGo/utils"
)

func TestNetworkLayerFunctionalities(t *testing.T) {
	// Initialize devices for two different local networks
	devices1 := []*physical.Device{
		{Name: "Device1", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 1},
		{Name: "Device2", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 2},
	}

	devices2 := []*physical.Device{
		{Name: "Device3", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 1},
		{Name: "Device4", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 2},
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
	router1Interface1 := &physical.Device{Name: "eth1", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 3}
	router1Interface2 := &physical.Device{Name: "eth2", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 4}
	router2Interface1 := &physical.Device{Name: "eth3", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 3}
	router2Interface2 := &physical.Device{Name: "eth4", LinkStatus: true, MACAddress: utils.GenerateRandomMAC(), PortNumber: 4}

	// Add and connect devices to the star topologies and connect switches to router interfaces
	starTopology1.AddEndDevice(devices1[0])
	starTopology1.AddEndDevice(devices1[1])
	starTopology1.AddEndDevice(router1Interface1)
	starTopology1.ConnectEndDevice(switchDevice1, router1Interface1)
	starTopology1.ConnectEndDevice(switchDevice1, devices1[0])
	starTopology1.ConnectEndDevice(switchDevice1, devices1[1])

	starTopology2.AddEndDevice(router2Interface1)
	starTopology2.AddEndDevice(devices2[0])
	starTopology2.AddEndDevice(devices2[1])
	starTopology2.ConnectEndDevice(switchDevice2, router2Interface1)
	starTopology2.ConnectEndDevice(switchDevice2, devices2[0])
	starTopology2.ConnectEndDevice(switchDevice2, devices2[1])

	for _, device := range starTopology1.EndDevices {
		t.Logf("Assigned MAC address %s to %s", device.MACAddress, device.Name)
	}
	t.Logf("Assigned MAC address %s to %s", router1Interface2.MACAddress, router1Interface2.Name)
	for _, device := range starTopology2.EndDevices {
		t.Logf("Assigned MAC address %s to %s", device.MACAddress, device.Name)
	}
	t.Logf("Assigned MAC address %s to %s", router2Interface2.MACAddress, router2Interface2.Name)

	utils.AssignIPAddress(starTopology1.EndDevices, "10.0.0.0/24")

	utils.AssignIPAddress(starTopology2.EndDevices, "20.0.0.0/24")
	router1Interface2.IPAddress = "30.0.0.1/24"
	router2Interface2.IPAddress = "30.0.0.2/24"

	for _, device := range starTopology1.EndDevices {
		t.Logf("Assigned IP address %s to %s", device.IPAddress, device.Name)
	}
	t.Logf("Assigned IP address %s to %s", router1Interface2.IPAddress, router1Interface2.Name)
	for _, device := range starTopology2.EndDevices {
		t.Logf("Assigned IP address %s to %s", device.IPAddress, device.Name)
	}
	t.Logf("Assigned IP address %s to %s", router2Interface2.IPAddress, router2Interface2.Name)

	for _, device := range starTopology1.EndDevices {
		device.ARPTable = make(map[string]physical.ARPEntry)
	}

	for _, device := range starTopology2.EndDevices {
		device.ARPTable = make(map[string]physical.ARPEntry)
	}

	// Initialize the routers and their ARP tables for testing
	router1 := network.Router{Interfaces: []physical.Device{*router1Interface1, *router1Interface2}, ID: "1"}

	router2 := network.Router{Interfaces: []physical.Device{*router2Interface1, *router2Interface2}, ID: "2"}

	// Connect the two routers
	router1.ConnectedTo = append(router1.ConnectedTo, &router2)
	router2.ConnectedTo = append(router2.ConnectedTo, &router1)

	//create static routing
	router1.AddRoutingEntry("10.0.0.0/24", "10.0.0.3/24", "eth1")
	router1.AddRoutingEntry("30.0.0.0/24", "30.0.0.1/24", "eth2")
	router1.AddRoutingEntry("20.0.0.0/24", "30.0.0.2/24", "eth4")

	router2.AddRoutingEntry("20.0.0.0/24", "20.0.0.1/24", "eth3")
	router2.AddRoutingEntry("10.0.0.0/24", "30.0.0.1/24", "eth2")
	router2.AddRoutingEntry("30.0.0.0/24", "30.0.0.2/24", "eth4")

	router1.DisplayRoutingTable()
	router2.DisplayRoutingTable()

	t.Logf("Network topology created successfully")

	sourceDevice := devices1[0]
	destinationDevice := devices2[1]
	gatewaydevice := starTopology1.Switch.EndDevices[2]
	packet := topologies.Packet{
		SourceIp:       sourceDevice.IPAddress,
		DestinationIp:  destinationDevice.IPAddress,
		SequenceNumber: 1,
		Acknowledgment: false,
		IsEnd:          false,
		Data:           []byte("Hello from Star1 to Star2"),
		Checksum:       utils.CalculateChecksum([]byte("Hello from Star1 to Star2")),
		Retransmit:     false,
	}

	ok, _ := utils.SameNid(sourceDevice.IPAddress, destinationDevice.IPAddress)

	if ok {

		if len(sourceDevice.ARPTable) == 0 || sourceDevice.ARPTable[destinationDevice.IPAddress] == (physical.ARPEntry{}) {
			// sourceDevice.CreateArpRequest(destinationDevice, &starTopology1.Switch)

			// starTopology1.SendDataToSwitch(&starTopology1.Switch, sourceDevice, receivingDevice, packet)
			// time.Sleep(3 * time.Second)
			LocalNetwork(sourceDevice, destinationDevice, starTopology1, packet)
		} else {
			// find it in cache

		}

	} else {
		// outside local network , just send it to gateway
		fmt.Printf("Device is outside local network , learning gateway macaddress\n")
		sourceDevice.CreateArpRequest(gatewaydevice, &starTopology1.Switch)
		gatewayMac := sourceDevice.ARPTable[gatewaydevice.IPAddress]
		gateDevice := utils.FindDeviceByMAC(starTopology1.EndDevices, gatewayMac.MACAddress)
		starTopology1.SendDataToSwitch(&starTopology1.Switch, sourceDevice, gateDevice, packet)
		time.Sleep(3 * time.Second)
		// nextHop := utils.LongestPrefixMatch(packet.DestinationIp, router1.RoutingTable)
		// fmt.Printf("interface is %s\n", nextHop)

		localDevice := router1.ForwardPacketToInterface(packet.DestinationIp)
		LocalNetwork(localDevice, destinationDevice, starTopology2, packet)
	}

	// Display ARP tables of devices
	for _, device := range starTopology1.EndDevices {
		t.Logf("ARP Table for %s:", device.Name)
		for ip, entry := range device.ARPTable {
			t.Logf("IP: %s, MAC Address: %s", ip, entry.MACAddress)
		}
	}

	for _, device := range starTopology2.EndDevices {
		t.Logf("ARP Table for %s:", device.Name)
		for ip, entry := range device.ARPTable {
			t.Logf("IP: %s, MAC Address: %s", ip, entry.MACAddress)
		}
	}
	t.Logf("Switching Table for %s:", starTopology1.Switch.Name)
	for mac, port := range starTopology1.Switch.SwitchingTable {
		t.Logf("%s\t\t%d\n", mac, port)
	}
	t.Logf("Switching Table for %s:", starTopology2.Switch.Name)
	for mac, port := range starTopology2.Switch.SwitchingTable {
		t.Logf("%s\t\t%d\n", mac, port)
	}

}

func LocalNetwork(sourceDevice *physical.Device, destinationDevice *physical.Device, s *datalink.StarTopologyWithSwitch, packet topologies.Packet) {
	sourceDevice.CreateArpRequest(destinationDevice, &s.Switch)

	destinationMacadd := sourceDevice.ARPTable[destinationDevice.IPAddress]
	receivingDevice := utils.FindDeviceByMAC(s.EndDevices, destinationMacadd.MACAddress)
	s.SendDataToSwitch(&s.Switch, sourceDevice, receivingDevice, packet)
	time.Sleep(3 * time.Second)
}
