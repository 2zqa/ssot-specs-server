package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type Kernel struct {
	ID      uint   `json:"-"`
	SpecsID uint   `json:"-"`
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

func ValidateKernel(v *validator.Validator, kernel *Kernel) {
	// No validation rules
}

func (k Kernel) String() string {
	return k.Name + delimiter +
		k.Version
}
