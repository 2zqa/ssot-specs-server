package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateSwapDevice(t *testing.T) {
	cases := []struct {
		name       string
		swapDevice SwapDevice
		expected   bool
	}{
		{
			name: "valid swap device",
			swapDevice: SwapDevice{
				Name: "swap",
				Size: 1024,
			},
			expected: true,
		},
		{
			name:       "empty swap device",
			swapDevice: SwapDevice{},
			expected:   true,
		},
		{
			name: "invalid size",
			swapDevice: SwapDevice{
				Size: -1,
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateSwapDevice(v, &tc.swapDevice)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
