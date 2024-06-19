// utils/ip_assigner.go
package utils

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/Mayvid0/netSimGo/internal/network"
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

func AssignIPAddress(devices interface{}, subnetCIDR string) error {
	var deviceList []*physical.Device

	switch v := devices.(type) {
	case []*physical.Device:
		deviceList = v
	case *physical.Device:
		deviceList = append(deviceList, v)
	default:
		return fmt.Errorf("unsupported type: %T", devices)
	}

	_, ipNet, err := net.ParseCIDR(subnetCIDR)
	if err != nil {
		return fmt.Errorf("invalid subnet CIDR: %v", err)
	}

	assignedIPs := make(map[string]bool)
	startIP := ipNet.IP.Mask(ipNet.Mask)
	incrementIP(startIP)

	for _, device := range deviceList {
		for ip := startIP; ipNet.Contains(ip); incrementIP(ip) {
			ipStr := ip.String()
			if !isIPAssigned(ipStr, assignedIPs) {
				device.IPAddress = fmt.Sprintf("%s/%d", ipStr, maskToCIDR(ipNet.Mask))
				assignedIPs[ipStr] = true
				break
			}
		}
	}

	return nil
}

// IncrementIP increments an IP address by one.
func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] != 0 {
			break
		}
	}
}

// MaskToCIDR converts a subnet mask to CIDR notation.
func maskToCIDR(mask net.IPMask) int {
	ones, _ := mask.Size()
	return ones
}

// IsIPAssigned checks if an IP address is already assigned.
func isIPAssigned(ip string, assignedIPs map[string]bool) bool {
	return assignedIPs[ip]
}

func SameNid(ip1 string, ip2 string) (bool, error) {
	ip1Addr, ip1Net, err := net.ParseCIDR(ip1)
	if err != nil {
		return false, fmt.Errorf("error parsing IP address %s: %v", ip1, err)
	}

	ip2Addr, ip2Net, err := net.ParseCIDR(ip2)
	if err != nil {
		return false, fmt.Errorf("error parsing IP address %s: %v", ip2, err)
	}

	// Check if the IP addresses are in the same network by comparing network addresses
	return ip1Addr.Mask(ip1Net.Mask).Equal(ip2Addr.Mask(ip2Net.Mask)), nil
}

func LongestPrefixMatch(ipWithMask string, routingTable []network.RoutingEntry) string {
	var longestMatch string

	_, ipNet, err := net.ParseCIDR(ipWithMask)
	if err != nil {
		fmt.Println("Error parsing IP with mask:", err)
		return longestMatch
	}

	maxPrefixLen := -1
	for _, entry := range routingTable {
		_, routeNet, err := net.ParseCIDR(entry.NetworkAddressMask)
		if err != nil {
			fmt.Println("Error parsing routing table entry:", err)
			continue
		}

		// Calculate prefix length manually
		ipPrefixLen, _ := routeNet.Mask.Size()
		if ipNet.Contains(routeNet.IP) && ipPrefixLen > maxPrefixLen {
			maxPrefixLen = ipPrefixLen
			longestMatch = entry.NextHop
		}
	}

	return longestMatch
}
