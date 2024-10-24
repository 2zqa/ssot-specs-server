package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"github.com/google/uuid"
)

var testingUUID = uuid.MustParse("fb249a22-f108-45ac-b30f-3baa78d09223")

func TestValidateDevice(t *testing.T) {
	cases := []struct {
		name     string
		device   Device
		expected bool
	}{
		{
			name: "valid cpu",
			device: Device{
				ID: testingUUID,
			},
			expected: true,
		},
		{
			name:     "no UUID",
			device:   Device{},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateDevice(v, &tc.device)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
