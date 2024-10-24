package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type NetworkInterfaceDriver struct {
	ID                 uint   `json:"-"`
	NetworkInterfaceID uint   `json:"-"`
	Name               string `json:"name,omitempty"`
	Version            string `json:"version,omitempty"`
	FirmwareVersion    string `json:"firmware_version,omitempty"`
}

func ValidateNetworkInterfaceDriver(v *validator.Validator, networkInterfaceDriver *NetworkInterfaceDriver) {
	// No validation rules
}

func (nid NetworkInterfaceDriver) String() string {
	return nid.Name + delimiter +
		nid.Version + delimiter +
		nid.FirmwareVersion
}
