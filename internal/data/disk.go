package data

import (
	"strconv"
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type Disk struct {
	ID            uint        `json:"-"`
	SpecsID       uint        `json:"-"`
	Name          string      `json:"name,omitempty"`
	SizeMegabytes int32       `json:"size_megabytes,omitempty"`
	Partitions    []Partition `json:"partitions,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

func ValidateDisk(v *validator.Validator, disks *Disk) {
	v.Check(disks.SizeMegabytes >= 0, "size_megabytes", greaterOrEqualZeroErrMsg)

	for _, partition := range disks.Partitions {
		ValidatePartition(v, &partition)
	}
}

func (d Disk) String() string {
	var partitions []string
	for _, partition := range d.Partitions {
		partitions = append(partitions, partition.String())
	}
	return d.Name + delimiter +
		strconv.Itoa(int(d.SizeMegabytes)) + delimiter +
		strings.Join(partitions, delimiter)
}
