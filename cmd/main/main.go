package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

func main() {
	var choice int
	fmt.Println("Select a topology:")
	fmt.Println("1. Point-to-Point")
	fmt.Println("2. Star")
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		pointToPointDriver()
	case 2:
		starDriver()
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
