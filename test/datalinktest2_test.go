package test

import (
	"fmt"
	"testing"

	"github.com/Mayvid0/netSimGo/internal/physical"
	"github.com/Mayvid0/netSimGo/internal/topologies"
	"github.com/Mayvid0/netSimGo/utils"
)

func TestNetworkCommunication(t *testing.T) {
	// Create two star topologies
	star1 := &topologies.Star{
		Hub: physical.Hub{
			Device: physical.Device{Name: "Hub1"},
		},
	}
	star2 := &topologies.Star{
		Hub: physical.Hub{
			Device: physical.Device{Name: "Hub2"},
		},
	}

	// Create five end devices for each star topology
	for i := 1; i <= 5; i++ {
		device1 := &physical.Device{
			Name:       fmt.Sprintf("Device%d_Star1", i),
			MACAddress: utils.GenerateRandomMAC(),
			LinkStatus: true,
			PortNumber: i,
			HasToken:   false,
		}
		device2 := &physical.Device{
			Name:       fmt.Sprintf("Device%d_Star2", i),
			MACAddress: utils.GenerateRandomMAC(),
			LinkStatus: true,
			PortNumber: i,
			HasToken:   false,
		}
		star1.AddEndDevice(device1)
		star2.AddEndDevice(device2)
	}

	// Create a bridge to connect the two hubs
	bridge := &physical.Bridge{
		Device: physical.Device{Name: "Bridge"},
	}
	bridge.EndDevices = append(bridge.EndDevices, &star1.Hub, &star2.Hub)

	packet := topologies.Packet{
		SequenceNumber: 1,
		Acknowledgment: false,
		IsEnd:          false,
		Data:           []byte("Hello from Star1 to Star2"),
		Checksum:       utils.CalculateChecksum([]byte("Hello from Star1 to Star2")),
		Retransmit:     false,
	}

	// Send packet from source to destination
	star1.SendDataToBridge(&star1.Hub, &star2.Hub, bridge, []topologies.Packet{packet})

}
