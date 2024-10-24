package data

import (
	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type OEM struct {
	ID           uint   `json:"-"`
	SpecsID      uint   `json:"-"`
	Manufacturer string `json:"manufacturer,omitempty"`
	ProductName  string `json:"product_name,omitempty"`
	SerialNumber string `json:"serial_number,omitempty"`
}

func ValidateOEM(v *validator.Validator, oem *OEM) {
	// No validation rules
}

func (o OEM) String() string {
	return o.Manufacturer + delimiter +
		o.ProductName + delimiter +
		o.SerialNumber
}
