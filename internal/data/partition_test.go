package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidatePartition(t *testing.T) {
	cases := []struct {
		name      string
		partition Partition
		expected  bool
	}{
		{
			name: "valid partition",
			partition: Partition{
				Filesystem:        "ext4",
				CapacityMegabytes: 1024,
				Source:            "/dev/sda1",
				Target:            "/",
			},
			expected: true,
		},
		{
			name:      "empty partition",
			partition: Partition{},
			expected:  true,
		},
		{
			name: "invalid capacity",
			partition: Partition{
				CapacityMegabytes: -1,
			},
			expected: false,
		},
		{
			name: "invalid target no slash",
			partition: Partition{
				Target: "boot/efi",
			},
			expected: false,
		},
		{
			name: "invalid target trailing slash",
			partition: Partition{
				Target: "/boot/efi/",
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidatePartition(v, &tc.partition)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
