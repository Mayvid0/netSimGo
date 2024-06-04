package datalink

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

type StarTopologyWithSwitch struct {
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

	log.Printf("Connection established between %s and %s\n", switchDevice.Name, device.Name)
	return nil
}

func (s *StarTopologyWithSwitch) SwitchingTable(switchDevice *physical.Switch) {
	for _, device := range s.EndDevices {
		switchDevice.SwitchingTable[device.MACAddress] = device.PortNumber
	}

	s.DisplaySwitchingTable(switchDevice)
}

func (s *StarTopologyWithSwitch) DisplaySwitchingTable(switchDevice *physical.Switch) {
	log.Println("Switching Table:")
	log.Println("MAC Address\tPort Number")
	for mac, port := range switchDevice.SwitchingTable {
		log.Printf("%s\t\t%d\n", mac, port)
	}
}

func (s *StarTopologyWithSwitch) SendDataFromSwitch(port int, packet topologies.Packet) {
	var targetDevice *physical.Device
	for _, device := range s.EndDevices {
		if device.PortNumber == port {
			targetDevice = device
			break
		}
	}

	if targetDevice != nil {
		log.Printf("Sending frame %d from Switch %s to %s\n", packet.SequenceNumber, s.Name, targetDevice.Name)
		s.ReceiveDataFromSwitch(targetDevice, packet)
	} else {
		log.Println("Error: End device not found for port number", port)
	}
}

func (s *StarTopologyWithSwitch) ReceiveDataFromSwitch(receiver *physical.Device, packet topologies.Packet) {
	log.Printf("Device %s received frame %d from Switch %s\n", receiver.Name, packet.SequenceNumber, s.Name)
	receivedString := string(packet.Data)
	log.Println("Received data:", receivedString)
}

func (s *StarTopologyWithSwitch) SendDataToSwitch(switchDevice *physical.Switch, source *physical.Device, receiver *physical.Device, packet topologies.Packet) {
	port, ok := switchDevice.SwitchingTable[receiver.MACAddress]
	if ok {
		log.Printf("Sending frame from %s to Switch %s\n", source.Name, switchDevice.Name)
		s.SendDataFromSwitch(port, packet)
	} else {
		log.Printf("Performing address learning\n")
		s.SwitchingTable(switchDevice)
		s.SendDataToSwitch(switchDevice, source, receiver, packet)
	}
}

func (s *StarTopologyWithSwitch) ReceiveAckFromSwitch(receiver *physical.Device, packet topologies.Packet) <-chan int {
	ackChannel := make(chan int)
	packet.Acknowledgment = true

	go func() {
		randomDelay := time.Duration(rand.Intn(5)) // Random delay up to 10 seconds
		time.Sleep(time.Second * randomDelay)
		for {
			if packet.Acknowledgment {
				ackChannel <- packet.SequenceNumber
			}
		}
	}()

	return ackChannel
}

func (s *StarTopologyWithSwitch) InitiateSelectiveRepeat(source *physical.Device, receiver *physical.Device, swi *physical.Switch, packets []topologies.Packet) {
	windowSize := 4
	rightPointer := 0
	leftPointer := 0
	receivedAcksLeft := make(map[int]bool)
	timers := make(map[int]*time.Timer)

	for rightPointer < len(packets) && leftPointer < len(packets) {
		for leftPointer <= rightPointer {
			// fmt.Printf("rightPointer: %d, len(packets): %d\n", rightPointer, len(packets))

			tokenAvailable := source.HasToken
			if !tokenAvailable {
				fmt.Println("Token not available, waiting...")
				time.Sleep(10 * time.Second)
				continue // Skip this iteration if token is not available
			}

			if rightPointer-leftPointer >= windowSize || rightPointer >= len(packets) {
				fmt.Println("All frames sent and acknowledged.")
				break
			}

			packetToSend := packets[rightPointer]
			s.SendDataToSwitch(swi, source, receiver, packetToSend)

			rightPointer++

			timer := time.AfterFunc(time.Second*3, func() {
				fmt.Printf("Timer expired and no ack was received for frame %d......Retransmitting....\n", packetToSend.SequenceNumber)
				s.RetransmitPacket(source, receiver, packetToSend, swi)
			})
			timers[packetToSend.SequenceNumber] = timer

			select {
			case ackPacketNumber := <-s.ReceiveAckFromSwitch(receiver, packetToSend):
				if ackPacketNumber >= 0 && ackPacketNumber < len(packets) {
					if !receivedAcksLeft[ackPacketNumber] {
						receivedAcksLeft[ackPacketNumber] = true
						log.Printf("Received ack of frame %d\n", ackPacketNumber)

						if timer, ok := timers[ackPacketNumber]; ok {
							// fmt.Printf("stopping timer of packet %d\n", ackPacketNumber)
							timer.Stop()
							delete(timers, ackPacketNumber)
						}

						for leftPointer < len(packets) && leftPointer <= rightPointer && receivedAcksLeft[packets[leftPointer].SequenceNumber] {
							delete(receivedAcksLeft, packets[leftPointer].SequenceNumber)
							leftPointer++
							// fmt.Printf("left pointer is %d\n", leftPointer)
						}

					}
				} else {
					fmt.Printf("Invalid acknowledgment frame number: %d\n", ackPacketNumber)
				}

			}
		}
	}

}

func (s *StarTopologyWithSwitch) RetransmitPacket(source *physical.Device, receiver *physical.Device, packetToSend topologies.Packet, swi *physical.Switch) {
	log.Printf("Re transmitting frame %d \n", packetToSend.SequenceNumber)
	s.SendDataToSwitch(swi, source, receiver, packetToSend)
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
			fmt.Printf("Token assigned to device %s\n", devices[nextIndex].Name)

			// Move to the next device
			currentIndex = nextIndex
		}
	}
}

// Divide the message into fixed size chunks and create packets
func CreatePackets(message string, maxPacketSize int) []topologies.Packet {
	var packets []topologies.Packet

	// Convert the message to bits
	messageBits := stringToBits(message)

	// Split the message bits into chunks
	chunks := splitBits(messageBits, maxPacketSize)

	// Calculate the total number of packets
	totalPackets := len(chunks)

	// Create packets with sequence numbers
	for i, chunk := range chunks {
		checksum := calculateChecksum(chunk) // Calculate checksum for each chunk
		packet := topologies.Packet{
			SequenceNumber: i % 8,
			Acknowledgment: false, // Data packet
			Data:           chunk,
			Checksum:       checksum, // Set checksum for the packet

		}

		// Check if it is the last chunk and set isEnd attribute
		if i == totalPackets-1 {
			packet.IsEnd = true
		}

		packets = append(packets, packet)
	}

	return packets
}

// Convert string to bits
func stringToBits(str string) []byte {
	var bits []byte
	for _, char := range str {
		bit := byte(char)
		bits = append(bits, bit)
	}
	return bits
}

// Split the message bits into fixed size chunks
func splitBits(bits []byte, chunkSize int) [][]byte {
	var chunks [][]byte
	for i := 0; i < len(bits); i += chunkSize {
		end := i + chunkSize
		if end > len(bits) {
			end = len(bits)
		}
		chunks = append(chunks, bits[i:end])
	}
	return chunks
}

// Calculate checksum for a byte array
func calculateChecksum(data []byte) uint16 {
	var sum uint32

	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}

	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}

	for sum>>16 > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	// Take the one's complement of the sum
	checksum := uint16(^sum)

	return checksum
}
