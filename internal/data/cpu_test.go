package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateCPU(t *testing.T) {
	cases := []struct {
		name     string
		cpu      CPU
		expected bool
	}{
		{
			name: "valid cpu",
			cpu: CPU{
				Architecture:          "x86_64",
				CoreCount:             4,
				CPUCount:              2,
				MaxFrequencyMegahertz: 3000,
			},
			expected: true,
		},
		{
			name:     "empty CPU",
			cpu:      CPU{},
			expected: true,
		},
		{
			name: "invalid architecture",
			cpu: CPU{
				Architecture: "invalid",
			},
			expected: false,
		},
		{
			name: "negative core count",
			cpu: CPU{
				CoreCount: -1,
			},
			expected: false,
		},
		{
			name: "negative cpu count",
			cpu: CPU{
				CPUCount: -1,
			},
			expected: false,
		},
		{
			name: "negative max frequency",
			cpu: CPU{
				MaxFrequencyMegahertz: -1,
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateCPU(v, &tc.cpu)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
