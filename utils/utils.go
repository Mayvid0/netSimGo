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
