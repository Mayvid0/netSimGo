package datalink

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

// we must create a struct like star topology , but with switch

type StarTopologyWithSwitch struct {

	// import the topology functions and the switch struct created in physical layer
	physical.Switch
	topologies.Topology
}

type Token struct {
	Available bool
}

func (s *StarTopologyWithSwitch) AddEndDevice(device *physical.Device) {
	s.EndDevices = append(s.EndDevices, device)
}

func (s *StarTopologyWithSwitch) ConnectEndDevice(switchDevice *physical.Switch, device *physical.Device) error {
	if !switchDevice.LinkStatus || !device.LinkStatus {
		return fmt.Errorf("cannot connect devices %s and %s: one or both devices are not linked", switchDevice.Name, device.Name)
	}

	fmt.Printf("Connection established between switch %s and end device %s\n", switchDevice.Name, device.Name)
	return nil
}

func (s *StarTopologyWithSwitch) SwitchingTable(switchDevice *physical.Switch) {
	// Create a table of the mapping of each port to its corresponding MAC address
	// Perform address learning
	for _, device := range s.EndDevices {
		switchDevice.SwitchingTable[device.MACAddress] = device.PortNumber
	}

	// Display the switching table
	s.DisplaySwitchingTable(switchDevice)
}

func (s *StarTopologyWithSwitch) DisplaySwitchingTable(switchDevice *physical.Switch) {
	fmt.Println("Switching Table:")
	fmt.Println("MAC Address\tPort Number")
	for mac, port := range switchDevice.SwitchingTable {
		fmt.Printf("%s\t\t%d\n", mac, port)
	}
}

func (s *StarTopologyWithSwitch) SendDataFromSwitch(port int, message string) {
	// Find the end device with the specified port number
	var targetDevice *physical.Device
	for _, device := range s.EndDevices {
		if device.PortNumber == port {
			targetDevice = device
			break
		}
	}

	if targetDevice != nil {
		if message[0] == '1' || message[0] == '0' {
			fmt.Printf("Sending message from Switch %s to %s before adding some noise: %s\n", s.Name, targetDevice.Name, message)
			noisyInputMessage := addNoise(message)
			fmt.Printf("Sending message from Switch %s to %s after adding some noise: %s\n", s.Name, targetDevice.Name, noisyInputMessage)
			s.ReceiveDataFromSwitch(targetDevice, noisyInputMessage)
		} else {
			fmt.Printf("Sending message from Switch %s to %s before adding some noise: %s\n", s.Name, targetDevice.Name, message)

			s.ReceiveDataFromSwitch(targetDevice, message)
		}

		// Send the message to the target device here
	} else {
		fmt.Println("Error: End device not found for port number", port)
	}
}

func (s *StarTopologyWithSwitch) ReceiveDataFromSwitch(receiver *physical.Device, message string) {
	if message[0] == '0' || message[0] == '1' {
		fmt.Printf(" Device %s Received message: before hamming error detection and correction %s \n", receiver.Name, message)
		correctedData := hammingDecoding(message)
		fmt.Printf(" Device %s Received message: after hamming error detection and correction %s \n", receiver.Name, correctedData)
	} else {
		fmt.Printf(" Device %s Received message:  %s \n", receiver.Name, message)
	}

}

func (s *StarTopologyWithSwitch) SendDataToSwitch(switchDevice *physical.Switch, source *physical.Device, receiver *physical.Device, message string) {
	if len(message) == 4 && (message[0] == '0' || message[0] == '1') {
		// Encode the message if it's a bit
		encodedMessage := hammingEncoding(message)
		port, ok := switchDevice.SwitchingTable[receiver.MACAddress]
		if ok {
			fmt.Printf("Sending encoded message from %s to Switch %s: %s\n", source.Name, switchDevice.Name, encodedMessage)
			s.SendDataFromSwitch(port, encodedMessage)
		} else {
			// Perform address learning
			fmt.Printf("Performing address learning\n")
			s.SwitchingTable(switchDevice)
			s.SendDataToSwitch(switchDevice, source, receiver, encodedMessage)
		}
	} else {
		// Send the original message if it's not a single bit
		port, ok := switchDevice.SwitchingTable[receiver.MACAddress]
		if ok {
			fmt.Printf("Sending message from %s to Switch %s: %s\n", source.Name, switchDevice.Name, message)
			s.SendDataFromSwitch(port, message)
		} else {
			// Perform address learning
			fmt.Printf("Performing address learning\n")
			s.SwitchingTable(switchDevice)
			s.SendDataToSwitch(switchDevice, source, receiver, message)
		}
	}
}

func hammingEncoding(inputMessage string) string {

	// Create a 7-bit Hamming code array
	hammingCode := []byte("0000000")

	// 	// 7 6 5 4 3 2 1  actual hamming code
	// 	// d d d p d p p
	// 	// 0 1 2 3 4 5 6  hamming code array
	//  // 1 1 0 0 1 1 0
	hammingCode[0] = inputMessage[0]
	hammingCode[1] = inputMessage[1]
	hammingCode[2] = inputMessage[2]
	hammingCode[4] = inputMessage[3]

	hammingCode[3] = calculateParity(hammingCode[0], hammingCode[1], hammingCode[2])
	hammingCode[5] = calculateParity(hammingCode[0], hammingCode[1], hammingCode[4])
	hammingCode[6] = calculateParity(hammingCode[4], hammingCode[2], hammingCode[0])

	encodedstring := string(hammingCode)
	fmt.Println(encodedstring)

	return string(encodedstring)
}
func calculateParity(b1, b2, b3 byte) byte {
	count := (b1 - '0') + (b2 - '0') + (b3 - '0')
	if count%2 == 0 {
		return '0'
	}
	return '1'
}

func addNoise(message string) string {
	// Create a new local random generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random index from 0 to 6
	index := r.Intn(7)

	// Convert the message to a byte slice to modify individual bits
	messageBytes := []byte(message)

	// Toggle the bit at the random index
	if messageBytes[index] == '0' {
		messageBytes[index] = '1'
	} else {
		messageBytes[index] = '0'
	}

	// Convert the modified byte slice back to a string
	noisyMessage := string(messageBytes)
	return noisyMessage
}

func hammingDecoding(receivedNoisyMessage string) string {
	// 7 6 5 4 3 2 1  actual hamming code
	// d d d p d p p
	// 0 1 2 3 4 5 6  hamming code array
	// 1 1 0 0 1 1 0
	messageBytes := []byte(receivedNoisyMessage)
	parityBits := []byte("000")

	parityBits[0] = checkParity(receivedNoisyMessage[3], receivedNoisyMessage[2], receivedNoisyMessage[1], receivedNoisyMessage[0])
	parityBits[1] = checkParity(receivedNoisyMessage[5], receivedNoisyMessage[4], receivedNoisyMessage[1], receivedNoisyMessage[0])
	parityBits[2] = checkParity(receivedNoisyMessage[6], receivedNoisyMessage[4], receivedNoisyMessage[2], receivedNoisyMessage[0])

	errorIndex := binaryToDecimal(parityBits) // Convert parity bits to decimal to get the error index

	if errorIndex > 0 {
		if messageBytes[7-errorIndex] == 0 {
			messageBytes[7-errorIndex] = 1
		} else {
			messageBytes[7-errorIndex] = 0
		}
	}

	decodedData := []byte("0000")
	decodedData[0] = messageBytes[0]
	decodedData[1] = messageBytes[1]
	decodedData[2] = messageBytes[2]
	decodedData[3] = messageBytes[4]

	return string(decodedData)
}

func checkParity(a, b, c, d byte) byte {
	count := (a - '0') + (b - '0') + (c - '0') + (d - '0')
	if count%2 == 0 {
		return '0'
	}
	return '1'
}

func binaryToDecimal(bits []byte) int {
	value := 0
	for i := len(bits) - 1; i >= 0; i-- {
		if bits[i] == '1' {
			value += 1 << (len(bits) - 1 - i)
		}
	}
	return value
}

func TokenPassing(devices []*physical.Device, token *Token, tokenHoldTime time.Duration) {
	// Assign the initial token to the first device
	devices[0].HasToken = true
	fmt.Printf("Token assigned to device %s\n", devices[0].Name)

	currentIndex := 0
	for {
		// Check if the current device has the token
		if devices[currentIndex].HasToken {
			// Simulate the device holding the token for some time
			time.Sleep(tokenHoldTime)

			// Pass the token to the next device
			devices[currentIndex].HasToken = false
			nextIndex := (currentIndex + 1) % len(devices)
			devices[nextIndex].HasToken = true

			// Move to the next device
			currentIndex = nextIndex
		}
	}
}
