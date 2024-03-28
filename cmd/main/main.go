package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	datalink "github.com/Mayvid0/netSimGo/internal/dataLinkLayer"
	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

func main() {
	var choice int
	fmt.Println("Select a topology:")
	fmt.Println("1. Point-to-Point")
	fmt.Println("2. Star")
	fmt.Println("3. Switch")
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		pointToPointDriver()
	case 2:
		starDriver()
	case 3:
		switchDriver()
	default:
		fmt.Println("Invalid choice. Exiting.")
		return
	}
}

func GenerateRandomMAC() string {
	const hexChars = "0123456789ABCDEF"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var sb strings.Builder
	for i := 0; i < 6; i++ {
		if i > 0 {
			sb.WriteString(":")
		}
		sb.WriteByte(hexChars[r.Intn(len(hexChars))])
		sb.WriteByte(hexChars[r.Intn(len(hexChars))])
	}
	return sb.String()
}

func pointToPointDriver() {
	topology := &topologies.PointToPoint{}

	var numDevices int
	fmt.Print("Enter the number of devices to add: ")
	fmt.Scanln(&numDevices)

	for i := 1; i <= numDevices; i++ {
		var name string
		fmt.Printf("Enter name for Device %d: ", i)
		fmt.Scanln(&name)

		device := &physical.Device{Name: name, MACAddress: GenerateRandomMAC(), LinkStatus: true}
		topology.AddEndDevice(device)
	}

	// Print devices with their names and MAC addresses for user reference
	fmt.Println("Devices in the topology:")
	for _, device := range topology.Devices {
		fmt.Printf("Name: %s, MAC: %s\n", device.Name, device.MACAddress)
	}

	var device1MAC, device2MAC string
	fmt.Println("Enter MAC addresses of two devices to connect:")
	fmt.Print("Device 1 MAC address: ")
	fmt.Scanln(&device1MAC)
	fmt.Print("Device 2 MAC address: ")
	fmt.Scanln(&device2MAC)

	// Find the devices in the topology based on their MAC addresses
	var dev1, dev2 *physical.Device
	for _, device := range topology.Devices {
		if device.MACAddress == device1MAC {
			dev1 = device
		} else if device.MACAddress == device2MAC {
			dev2 = device
		}
	}

	if dev1 == nil || dev2 == nil {
		fmt.Println("Error: Devices not found in topology.")
		return
	}

	err := topology.ConnectEndDevice(dev1, dev2)
	if err != nil {
		fmt.Println("Error connecting devices:", err)
		return
	}

	fmt.Print("Enter message to send: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	message = message[:len(message)-1]

	topology.SendData(dev1, dev2, message)
}

func starDriver() {
	topology := &topologies.Star{}

	var hubName string
	fmt.Print("Enter name for Hub: ")
	fmt.Scanln(&hubName)

	hub := &physical.Hub{Device: physical.Device{Name: hubName, MACAddress: GenerateRandomMAC(), LinkStatus: true}}

	var numDevices int
	fmt.Print("Enter the number of devices to add: ")
	fmt.Scanln(&numDevices)

	for i := 1; i <= numDevices; i++ {
		var name string
		fmt.Printf("Enter name for Device %d: ", i)
		fmt.Scanln(&name)

		device := &physical.Device{Name: name, MACAddress: GenerateRandomMAC(), LinkStatus: true}
		topology.AddEndDevice(device)
		topology.ConnectEndDevice(hub, device)
	}

	// Print devices with their names and MAC addresses for user reference
	fmt.Println("Devices in the topology:")
	fmt.Printf("Hub: Name: %s, MAC: %s\n", hubName, hub.MACAddress)
	for _, device := range topology.EndDevices {
		fmt.Printf("Name: %s, MAC: %s\n", device.Name, device.MACAddress)
	}

	fmt.Print("Enter device that has to send message (by MAC address): ")
	var sourceMAC string
	fmt.Scanln(&sourceMAC)

	var source *physical.Device
	for _, device := range topology.EndDevices {
		if device.MACAddress == sourceMAC {
			source = device
			break
		}
	}
	if source == nil {
		fmt.Println("Error: Device not found.")
		return
	}

	fmt.Print("Enter message to send: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	message = message[:len(message)-1]

	topology.SendDataToHub(source, hub, message)
}

func switchDriver() {
	switchTopology := &datalink.StarTopologyWithSwitch{}

	// Ask user for switch details
	var switchName string
	fmt.Print("Enter name for Switch: ")
	fmt.Scanln(&switchName)

	var totalPorts int
	fmt.Print("Enter total number of ports for the switch: ")
	fmt.Scanln(&totalPorts)

	switchDevice := &physical.Switch{
		Hub: physical.Hub{
			Device: physical.Device{
				Name:       switchName,
				MACAddress: GenerateRandomMAC(),
				LinkStatus: true,
			},
			NumberPorts: totalPorts,
		},
		SwitchingTable: make(map[string]int),
	}

	switchTopology.Switch = *switchDevice

	var numEndDevices int
	fmt.Print("Enter the number of end devices to add: ")
	fmt.Scanln(&numEndDevices)
	switchTopology.Switch.PortNumber = 1

	// Add end devices and connect them to the switch
	for i := 1; i <= totalPorts; i++ {
		var deviceName string
		fmt.Printf("Enter name for End Device %d: ", i)
		fmt.Scanln(&deviceName)

		deviceMAC := GenerateRandomMAC()
		devicePort := switchTopology.Switch.PortNumber // Assign port number
		endDevice := &physical.Device{
			Name:       deviceName,
			MACAddress: deviceMAC,
			LinkStatus: true,
			PortNumber: devicePort,
		}

		// Connect end device to switch
		switchTopology.AddEndDevice(endDevice)
		switchTopology.ConnectEndDevice(&switchTopology.Switch, endDevice)

		// Update port number for the next device
		switchTopology.Switch.PortNumber++
	}

	// Display all devices with their MAC addresses
	fmt.Println("Devices in the switch topology:")
	for _, device := range switchTopology.EndDevices {
		fmt.Printf("Name: %s, MAC: %s\n", device.Name, device.MACAddress)
	}

	// Ask user for source and receiver devices
	fmt.Print("Enter MAC address of the source device: ")
	var sourceMAC string
	fmt.Scanln(&sourceMAC)

	fmt.Print("Enter MAC address of the receiver device: ")
	var receiverMAC string
	fmt.Scanln(&receiverMAC)

	// Find the source and receiver devices
	var source *physical.Device
	var receiver *physical.Device
	for _, device := range switchTopology.EndDevices {
		if device.MACAddress == sourceMAC {
			source = device
		} else if device.MACAddress == receiverMAC {
			receiver = device
		}
	}

	if source == nil || receiver == nil {
		fmt.Println("Error: Source or receiver device not found.")
		return
	}

	// Send message from source to receiver
	fmt.Print("Enter message to send: ")
	reader := bufio.NewReader(os.Stdin)
	message, _ := reader.ReadString('\n')
	message = message[:len(message)-1]

	switchTopology.SendDataToSwitch(&switchTopology.Switch, source, receiver, message)
}
