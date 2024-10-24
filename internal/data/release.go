package data

import "github.com/2zqa/ssot-specs-server/internal/validator"

type Release struct {
	ID       uint   `json:"-"`
	SpecsID  uint   `json:"-"`
	Name     string `json:"name,omitempty"`
	Version  string `json:"version,omitempty"`
	Codename string `json:"codename,omitempty"`
}

func ValidateRelease(v *validator.Validator, release *Release) {
	// No validation rules
}

func (r Release) String() string {
	return r.Name + delimiter +
		r.Version + delimiter +
		r.Codename
}
