package data

import (
	"strconv"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

type Partition struct {
	ID                uint   `json:"-"`
	DiskID            uint   `json:"-"`
	Filesystem        string `json:"filesystem,omitempty"`
	CapacityMegabytes int32  `json:"capacity_megabytes,omitempty"`
	Source            string `json:"source,omitempty"`
	Target            string `json:"target,omitempty"`
}

func ValidatePartition(v *validator.Validator, partition *Partition) {
	v.Check(partition.CapacityMegabytes >= 0, "capacity_megabytes", greaterOrEqualZeroErrMsg)
	// If Target is not empty, check if it is a valid path
	if partition.Target != "" {
		v.Check(partition.Target[0] == '/', "target", "paths must start with a slash")
		if len(partition.Target) > 1 {
			v.Check(partition.Target[len(partition.Target)-1] != '/', "target", "paths must not end with a slash")
		}
	}
}

func (p Partition) String() string {
	return p.Filesystem + delimiter +
		strconv.Itoa(int(p.CapacityMegabytes)) + delimiter +
		p.Source + delimiter +
		p.Target
}
