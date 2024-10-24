package data

import (
	"strconv"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type SwapDevice struct {
	ID       uint   `json:"-"`
	MemoryID uint   `json:"-"`
	Name     string `json:"name,omitempty"`
	Size     int32  `json:"size,omitempty"`
}

func ValidateSwapDevice(v *validator.Validator, swapDevice *SwapDevice) {
	v.Check(swapDevice.Size >= 0, "size", greaterOrEqualZeroErrMsg)
}

func (s SwapDevice) String() string {
	return s.Name + delimiter +
		strconv.Itoa(int(s.Size))
}
