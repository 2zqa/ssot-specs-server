package data

import (
	"strconv"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type DIMM struct {
	ID            uint   `json:"-"`
	SpecsID       uint   `json:"-"`
	SizeGigabytes int32  `json:"size_gigabytes,omitempty"`
	Speed         int32  `json:"speed_mt_s,omitempty"`
	Manufacturer  string `json:"manufacturer,omitempty"`
	SerialNumber  string `json:"serial_number,omitempty"`
	Type          string `json:"type,omitempty"`
	PartNumber    string `json:"part_number,omitempty"`
	FormFactor    string `json:"form_factor,omitempty"`
	Locator       string `json:"locator,omitempty"`
	BankLocator   string `json:"bank_locator,omitempty"`
}

func ValidateDIMM(v *validator.Validator, dimm *DIMM) {
	v.Check(dimm.SizeGigabytes >= 0, "size_gigabytes", greaterOrEqualZeroErrMsg)
	v.Check(dimm.Speed >= 0, "speed_mt_s", greaterOrEqualZeroErrMsg)
}

func (d DIMM) String() string {
	return strconv.Itoa(int(d.SizeGigabytes)) + delimiter +
		strconv.Itoa(int(d.Speed)) + delimiter +
		d.Manufacturer + delimiter +
		d.SerialNumber + delimiter +
		d.Type + delimiter +
		d.PartNumber + delimiter +
		d.FormFactor + delimiter +
		d.Locator + delimiter +
		d.BankLocator
}
