package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
)

func FindDeviceByMAC(devices []*physical.Device, macAddress string) *physical.Device {
	for _, device := range devices {
		if device.MACAddress == macAddress {
			return device
		}
	}
	return nil
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
func CalculateChecksum(data []byte) uint16 {
	var sum uint32

	// Sum all 16-bit words in the data
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 + uint32(data[i+1])
	}

	// If the data length is odd, add the last byte as a 16-bit word
	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}

	// Fold the 32-bit sum to 16 bits
	for sum>>16 > 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}

	// Take the one's complement of the sum
	checksum := uint16(^sum)

	return checksum
}
