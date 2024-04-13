package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	datalink "github.com/Mayvid0/netSimGo/internal/dataLinkLayer"
	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
	"github.com/Mayvid0/netSimGo/utils"
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
		fmt.Println("3. Hub")
		fmt.Println("4. Switch")
		var choice2 int
		fmt.Scanln(&choice2)
		switch choice2 {
		case 3:
			starDriver()
		case 4:
			switchDriver()
		}
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
		var name string
		fmt.Printf("Enter name for Device %d: ", i)
		fmt.Scanln(&name)

		device := &physical.Device{Name: name, MACAddress: utils.GenerateRandomMAC(), LinkStatus: true}
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

	hub := &physical.Hub{Device: physical.Device{Name: hubName, MACAddress: utils.GenerateRandomMAC(), LinkStatus: true}}

	var numDevices int
	fmt.Print("Enter the number of devices to add: ")
	fmt.Scanln(&numDevices)

	for i := 1; i <= numDevices; i++ {
		var name string
		fmt.Printf("Enter name for Device %d: ", i)
		fmt.Scanln(&name)

		device := &physical.Device{Name: name, MACAddress: utils.GenerateRandomMAC(), LinkStatus: true}
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

type MACPair struct {
	SourceMAC   string
	Message     string
	ReceiverMAC string
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
				MACAddress: utils.GenerateRandomMAC(),
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

		deviceMAC := utils.GenerateRandomMAC()
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

	//access control protocol
	token := &datalink.Token{Available: true}
	tokenChannel := make(chan bool)
	go func() {
		datalink.TokenPassing(switchTopology.EndDevices, token, 20*time.Second)
		tokenChannel <- true
	}()

	// Start the goroutine for prompting MAC addresses
	macChannel := make(chan MACPair)
	inputTrigger := make(chan bool, 1)
	inputTrigger <- true
	go func() {
		for {
			<-inputTrigger // Wait for trigger signal

			fmt.Println("Devices in the switch topology:")
			for _, device := range switchTopology.EndDevices {
				fmt.Printf("Name: %s, MAC: %s\n", device.Name, device.MACAddress)
			}
			fmt.Print("Enter MAC address of the source device (or 'exit' to stop): ")
			var sourceMAC string
			fmt.Scanln(&sourceMAC)

			if sourceMAC == "exit" {
				break
			}

			fmt.Print("Enter message to send : ")
			reader := bufio.NewReader(os.Stdin)
			message, _ := reader.ReadString('\n')
			message = message[:len(message)-1]

			fmt.Print("Enter MAC address of the receiver device: ")
			var receiverMAC string
			fmt.Scanln(&receiverMAC)

			macPair := MACPair{SourceMAC: sourceMAC, Message: message, ReceiverMAC: receiverMAC}
			macChannel <- macPair
		}
	}()

	// Main loop to handle token-based data transmission
	for {
		select {
		case macPair := <-macChannel:
			// Handle MAC pair input
			sourceDevice := utils.FindDeviceByMAC(switchTopology.EndDevices, macPair.SourceMAC)
			receiverDevice := utils.FindDeviceByMAC(switchTopology.EndDevices, macPair.ReceiverMAC)

			if sourceDevice == nil || receiverDevice == nil {
				fmt.Println("Error: Source or receiver device not found. Try again")
				inputTrigger <- true
				continue
			}

			if sourceDevice.HasToken {
				switchTopology.SendDataToSwitch(&switchTopology.Switch, sourceDevice, receiverDevice, macPair.Message)
				inputTrigger <- true
			} else {
				fmt.Println("Error: Source device does not have the token.")
				inputTrigger <- true
			}

		case <-tokenChannel:
			return
		}
	}

}
