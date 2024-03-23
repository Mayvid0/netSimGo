// topologies/topology.go
package topologies

import "github.com/Mayvid0/netSimGo/internal/physical"

type Topology interface {
	AddEndDevice(device *physical.Device)
	ConnectEndDevice(device1 *physical.Device, device2 *physical.Device) error
	SendData(source *physical.Device, receiver *physical.Device, message string)
	ReceiveData(source *physical.Device, receiver *physical.Device, message string)
}
