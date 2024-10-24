package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateDisk(t *testing.T) {
	cases := []struct {
		name     string
		disk     Disk
		expected bool
	}{
		{
			name: "valid disk",
			disk: Disk{
				SizeMegabytes: 100,
			},
			expected: true,
		},
		{
			name:     "empty disk",
			disk:     Disk{},
			expected: true,
		},
		{
			name: "invalid size",
			disk: Disk{
				SizeMegabytes: -1,
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateDisk(v, &tc.disk)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
