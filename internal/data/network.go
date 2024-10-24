package data

import (
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type Network struct {
	ID         uint               `json:"-"`
	SpecsID    uint               `json:"-"`
	Hostname   string             `json:"hostname,omitempty"`
	Interfaces []NetworkInterface `json:"interfaces,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

func ValidateNetwork(v *validator.Validator, network *Network) {
	for _, networkInterface := range network.Interfaces {
		ValidateNetworkInterface(v, &networkInterface)
	}
}

func (n Network) String() string {
	var interfaces []string
	for _, networkInterface := range n.Interfaces {
		interfaces = append(interfaces, networkInterface.String())
	}
	return n.Hostname + delimiter + strings.Join(interfaces, delimiter)
}
