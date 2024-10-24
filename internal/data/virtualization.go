package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type Virtualization struct {
	ID      uint   `json:"-"`
	SpecsID uint   `json:"-"`
	Type    string `json:"type,omitempty"`
}

func ValidateVirtualization(v *validator.Validator, virtualization *Virtualization) {
	// No validation rules
}

func (v Virtualization) String() string {
	return v.Type
}
