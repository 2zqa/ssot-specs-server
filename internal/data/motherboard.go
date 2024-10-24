package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type Motherboard struct {
	ID           uint   `json:"-"`
	SpecsID      uint   `json:"-"`
	Vendor       string `json:"vendor,omitempty"`
	Name         string `json:"name,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
}

func ValidateMotherboard(v *validator.Validator, motherboard *Motherboard) {
	// No validation rules
}

func (m Motherboard) String() string {
	return m.Vendor + delimiter +
		m.Name + delimiter +
		m.SerialNumber
}
