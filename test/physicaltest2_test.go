package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
)

func GenerateRandomMAC() string {
	const hexChars = "0123456789ABCDEF"
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	var mac string
	for i := 0; i < 6; i++ {
		if i > 0 {
			mac += ":"
		}
		mac += fmt.Sprintf("%c%c", hexChars[r.Intn(len(hexChars))], hexChars[r.Intn(len(hexChars))])
	}
	return mac
}

func TestStarTopology(t *testing.T) {
	star := &topologies.Star{}
	hub := &physical.Hub{Device: physical.Device{Name: "Centre hub", MACAddress: "00:11:22:33:44:55", LinkStatus: true}}

	// Adding end devices to the star topology
	for i := 1; i <= 5; i++ {
		deviceName := fmt.Sprintf("Device%d", i)
		deviceMAC := GenerateRandomMAC()
		device := &physical.Device{Name: deviceName, MACAddress: deviceMAC, LinkStatus: true}
		star.AddEndDevice(device)
	}

	// Connecting end devices to the hub
	for _, device := range star.EndDevices {
		err := star.ConnectEndDevice(hub, device)
		if err != nil {
			t.Errorf("Error connecting device %s to hub: %v", device.Name, err)
		}
	}

	// Sending data from one end device to the hub
	sourceDevice := star.EndDevices[0]
	message := "Hello, Hub, this is a unit test!"
	star.SendDataToHub(sourceDevice, hub, message)

}
