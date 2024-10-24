package data

import (
	"strings"
	"time"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/google/uuid"
)

type Specs struct {
	ID             uint           `json:"-"`
	DeviceID       uuid.UUID      `json:"-"`
	BootTime       time.Time      `json:"boot_time,omitempty"`
	Motherboard    Motherboard    `json:"motherboard,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	CPU            CPU            `json:"cpu,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Disks          []Disk         `json:"disks,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Network        Network        `json:"network,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Bios           Bios           `json:"bios,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Memory         Memory         `json:"memory,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Kernel         Kernel         `json:"kernel,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Release        Release        `json:"release,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	DIMMs          []DIMM         `json:"dimms,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	Virtualization Virtualization `json:"virtualization,omitempty" gorm:"constraint:OnDelete:CASCADE"`
	OEM            OEM            `json:"oem,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

func ValidateSpecs(v *validator.Validator, specs *Specs) {
	ValidateMotherboard(v, &specs.Motherboard)
	ValidateCPU(v, &specs.CPU)
	for _, disk := range specs.Disks {
		ValidateDisk(v, &disk)
	}
	ValidateNetwork(v, &specs.Network)
	ValidateBios(v, &specs.Bios)
	ValidateMemory(v, &specs.Memory)
	ValidateKernel(v, &specs.Kernel)
	ValidateRelease(v, &specs.Release)
	for _, dimm := range specs.DIMMs {
		ValidateDIMM(v, &dimm)
	}
	ValidateVirtualization(v, &specs.Virtualization)
	ValidateOEM(v, &specs.OEM)
}

func (s Specs) String() string {
	var disks []string
	for _, disk := range s.Disks {
		disks = append(disks, disk.String())
	}
	var dimms []string
	for _, dimm := range s.DIMMs {
		dimms = append(dimms, dimm.String())
	}

	return s.Motherboard.String() + delimiter +
		s.CPU.String() + delimiter +
		strings.Join(disks, delimiter) + delimiter +
		s.Network.String() + delimiter +
		s.Bios.String() + delimiter +
		s.Memory.String() + delimiter +
		s.Kernel.String() + delimiter +
		s.Release.String() + delimiter +
		strings.Join(dimms, delimiter) + delimiter +
		s.Virtualization.String() + delimiter +
		s.OEM.String()
}
