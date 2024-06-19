package network

import (
	"fmt"
	"net"
	"strings"

	"github.com/Mayvid0/netSimGo/internal/physical"
)

type RoutingEntry struct {
	NetworkAddressMask string
	NextHop            string
	Interface          string
}

type Router struct {
	ID              string
	Interfaces      []physical.Device
	RoutingTable    []RoutingEntry
	RoutingProtocol string
	ConnectedTo     []*Router
}

// AddRoutingEntry adds a new routing entry to the router's routing table
func (r *Router) AddRoutingEntry(networkAddressMask string, nextHop string, intf string) {
	entry := RoutingEntry{
		NetworkAddressMask: networkAddressMask,
		NextHop:            nextHop,
		Interface:          intf,
	}
	r.RoutingTable = append(r.RoutingTable, entry)
}

func (r *Router) DisplayRoutingTable() {
	fmt.Printf("Routing Table for Router %s:\n", r.ID)
	fmt.Println("Network Address/Mask\tNext Hop\tInterface")
	for _, entry := range r.RoutingTable {
		fmt.Printf("%s\t%s\t%s\n", entry.NetworkAddressMask, entry.NextHop, entry.Interface)
	}
}

func (r *Router) findLocalInterface(ip net.IP) (bool, *physical.Device) {
	for _, intf := range r.Interfaces {
		_, network, err := net.ParseCIDR(intf.IPAddress)
		if err != nil {
			continue
		}
		if network.Contains(ip) {
			return true, &intf
		}
	}
	return false, nil
}
func (r *Router) ForwardPacketToInterface(destinationIP string) *physical.Device {
	// Find the appropriate next hop based on the destination IP address
	var nextHop, outInterface string

	ip, _, _ := net.ParseCIDR(strings.TrimSpace(destinationIP))
	if ip == nil {
		fmt.Printf("Error parsing destination IP %s\n", destinationIP)
		return nil
	}

	// Check if the destination IP is within the local network
	if isLocal, intf := r.findLocalInterface(ip); isLocal {
		fmt.Printf("Destination IP %s is within the local network of router %s on interface %s. Stopping forwarding.\n", destinationIP, r.ID, intf.Name)
		return intf
	}

	for _, entry := range r.RoutingTable {
		trimmedNetworkAddressMask := strings.TrimSpace(entry.NetworkAddressMask)
		_, network, err := net.ParseCIDR(trimmedNetworkAddressMask)
		if err != nil {
			fmt.Printf("Error parsing network address %s: %v\n", trimmedNetworkAddressMask, err)
			continue
		}
		if network.Contains(ip) {
			nextHop = entry.NextHop
			outInterface = entry.Interface
			break
		}
	}

	if nextHop != "" {
		fmt.Printf("Forwarding packet from router %s to next hop %s via interface %s\n", r.ID, nextHop, outInterface)

		// Find the router associated with the next hop IP address
		for _, router := range r.ConnectedTo {
			for _, intf := range router.Interfaces {
				if intf.IPAddress == nextHop {
					fmt.Printf("Found next hop router %s, forwarding packet\n", router.ID)
					return router.ForwardPacketToInterface(destinationIP)
				}
			}
		}
		fmt.Printf("Next hop router with IP %s not found among connected routers\n", nextHop)
	} else {
		fmt.Println("No route found for packet forwarding")
	}

	return nil // Return nil if no interface is found or no route is found
}
