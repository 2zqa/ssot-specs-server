package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type Bios struct {
	ID      uint   `json:"-"`
	SpecsID uint   `json:"-"`
	Vendor  string `json:"vendor,omitempty"`
	Version string `json:"version,omitempty"`
	Date    string `json:"date,omitempty"`
}

func ValidateBios(v *validator.Validator, bios *Bios) {
	// Date nor version can be parsed because it is not in a guaranteed format
}

func (b Bios) String() string {
	return b.Vendor + delimiter + b.Version + delimiter + b.Date
}
