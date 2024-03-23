package main

import (
	"bufio"
	"fmt"
	"os"

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

func pointToPointDriver() {
	topology := &topologies.PointToPoint{}

	var numDevices int
	fmt.Print("Enter the number of devices to add: ")
	fmt.Scanln(&numDevices)

	for i := 1; i <= numDevices; i++ {
		var name, mac string
		fmt.Printf("Enter name for Device %d: ", i)
		fmt.Scanln(&name)
		fmt.Printf("Enter MAC Address for Device %d: ", i)
		fmt.Scanln(&mac)

		device := &physical.Device{Name: name, MACAddress: mac, LinkStatus: true}
		topology.AddEndDevice(device)
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
	// Driver code for Star topology
	// Add your logic here
}
