package data

import (
	"strconv"
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/lib/pq"
)

var (
	architectureList = []string{"x86_64", "i386", "i686", "armv6l", "armv7l", "aarch64", "ppc64le", "s390x"}
)

type CPU struct {
	ID                    uint           `json:"-"`
	SpecsID               uint           `json:"-"`
	Name                  string         `json:"name,omitempty"`
	Architecture          string         `json:"architecture,omitempty"`
	CoreCount             int32          `json:"core_count,omitempty"`
	CPUCount              int32          `json:"cpu_count,omitempty"`
	MaxFrequencyMegahertz int32          `json:"max_frequency_megahertz,omitempty"`
	Mitigations           pq.StringArray `json:"mitigations,omitempty" gorm:"type:text[]"`
}

func ValidateCPU(v *validator.Validator, cpu *CPU) {
	v.Check(cpu.Architecture == "" || validator.PermittedValue(cpu.Architecture, architectureList...), "architecture", "invalid architecture value")
	v.Check(cpu.CoreCount >= 0, "core_count", greaterOrEqualZeroErrMsg)
	v.Check(cpu.CPUCount >= 0, "cpu_count", greaterOrEqualZeroErrMsg)
	v.Check(cpu.MaxFrequencyMegahertz >= 0, "max_frequency_megahertz", greaterOrEqualZeroErrMsg)
}

func (c CPU) String() string {
	return c.Name + delimiter +
		c.Architecture + delimiter +
		strconv.Itoa(int(c.CoreCount)) + delimiter +
		strconv.Itoa(int(c.CPUCount)) + delimiter +
		strconv.Itoa(int(c.MaxFrequencyMegahertz)) + delimiter +
		strings.Join(c.Mitigations, delimiter)
}
