package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateDIMM(t *testing.T) {
	cases := []struct {
		name     string
		dimm     DIMM
		expected bool
	}{
		{
			name: "valid dimm",
			dimm: DIMM{
				SizeGigabytes: 1,
				Speed:         1,
			},
			expected: true,
		},
		{
			name:     "empty dimm",
			dimm:     DIMM{},
			expected: true,
		},
		{
			name: "invalid speed",
			dimm: DIMM{
				Speed: -1,
			},
			expected: false,
		},
		{
			name: "invalid size",
			dimm: DIMM{
				SizeGigabytes: -1,
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateDIMM(v, &tc.dimm)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
