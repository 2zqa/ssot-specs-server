package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateNetworkInterface(t *testing.T) {
	cases := []struct {
		name     string
		netIface NetworkInterface
		expected bool
	}{
		{
			name: "valid network interface",
			netIface: NetworkInterface{
				MACAddress: "ac:1f:6b:1a:7e:2e",
				IPv4Addresses: []string{
					// IP address reserved for documentation: https://datatracker.ietf.org/doc/html/rfc5737#section-3
					"198.51.100.0/24",
				},
				IPv6Addresses: []string{
					// IP address reserved for documentation: https://datatracker.ietf.org/doc/html/rfc3849
					"2001:db8:ffff:ffff:ffff:ffff:ffff:ffff/32",
				},
			},
			expected: true,
		},
		{
			name:     "empty network interface",
			netIface: NetworkInterface{},
			expected: true,
		},
		{
			name: "invalid ipv4 address",
			netIface: NetworkInterface{
				IPv4Addresses: []string{
					"192.168.1783.26/24",
				},
			},
			expected: false,
		},
		{
			name: "invalid ipv6 address",
			netIface: NetworkInterface{
				IPv6Addresses: []string{
					"fe80::ac28:7ae:9a8f:5f69/64",                // valid
					"2001:1c01:2e043:5000:63b:dc32:dfe6:35c4/64", // invalid
				},
			},
			expected: false,
		},
		{
			name: "invalid mac address",
			netIface: NetworkInterface{
				MACAddress: "ac:1f:6b:1a:7e:2e:00",
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateNetworkInterface(v, &tc.netIface)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
