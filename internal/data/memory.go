package data

import (
	"strconv"
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type Memory struct {
	ID          uint         `json:"-"`
	SpecsID     uint         `json:"-"`
	Memory      int32        `json:"memory,omitempty"`
	Swap        int32        `json:"swap,omitempty"`
	SwapDevices []SwapDevice `json:"swap_devices,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

func ValidateMemory(v *validator.Validator, memory *Memory) {
	v.Check(memory.Memory >= 0, "memory", greaterOrEqualZeroErrMsg)
	v.Check(memory.Swap >= 0, "swap", greaterOrEqualZeroErrMsg)
	for _, swapDevice := range memory.SwapDevices {
		ValidateSwapDevice(v, &swapDevice)
	}
}

func (m Memory) String() string {
	var swapDevices []string
	for _, swapDevice := range m.SwapDevices {
		swapDevices = append(swapDevices, swapDevice.String())
	}

	return strconv.Itoa(int(m.Memory)) + delimiter +
		strconv.Itoa(int(m.Swap)) + delimiter +
		strings.Join(swapDevices, delimiter)
}
