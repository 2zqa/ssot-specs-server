package data

import (
	"net"
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/lib/pq"
)

type NetworkInterface struct {
	ID            uint                   `json:"-"`
	NetworkID     uint                   `json:"-"`
	MACAddress    string                 `json:"mac_address,omitempty"`
	Driver        NetworkInterfaceDriver `json:"driver,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	IPv4Addresses pq.StringArray         `json:"ipv4_addresses,omitempty" gorm:"type:text[]"`
	IPv6Addresses pq.StringArray         `json:"ipv6_addresses,omitempty" gorm:"type:text[]"`
}

func ValidateNetworkInterface(v *validator.Validator, networkInterface *NetworkInterface) {
	ValidateNetworkInterfaceDriver(v, &networkInterface.Driver)

	if networkInterface.MACAddress != "" {
		_, err := net.ParseMAC(networkInterface.MACAddress)
		v.Check(err == nil, "mac_address", "must be a valid MAC address")
	}
	for _, ipv4Address := range networkInterface.IPv4Addresses {
		v.Check(isValidIPv4Address(ipv4Address), "ipv4_addresses", "must be a valid IPv4 address")
	}
	for _, ipv6Address := range networkInterface.IPv6Addresses {
		v.Check(isValidIPv6Address(ipv6Address), "ipv6_addresses", "must be a valid IPv6 address")
	}
}

func parseIP(ipString string) net.IP {
	ip, _, _ := net.ParseCIDR(ipString)
	return ip
}

func isValidIPv4Address(ipString string) bool {
	return parseIP(ipString).To4() != nil
}

func isValidIPv6Address(ipString string) bool {
	return parseIP(ipString).To16() != nil
}

func (n NetworkInterface) String() string {
	return n.MACAddress + delimiter +
		n.Driver.String() + delimiter +
		strings.Join(n.IPv4Addresses, delimiter) + delimiter +
		strings.Join(n.IPv6Addresses, delimiter)
}
