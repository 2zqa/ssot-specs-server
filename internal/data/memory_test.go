// BEGIN: 9f3c8a7b6d5c
package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateMemory(t *testing.T) {
	cases := []struct {
		name     string
		memory   Memory
		expected bool
	}{
		{
			name: "valid memory",
			memory: Memory{
				Memory: 1000000000,
				Swap:   1000000000,
			},
			expected: true,
		},
		{
			name:     "empty memory",
			memory:   Memory{},
			expected: true,
		},
		{
			name: "negative total memory",
			memory: Memory{
				Memory: -1,
			},
			expected: false,
		},
		{
			name: "negative total swap",
			memory: Memory{
				Swap: -1,
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateMemory(v, &tc.memory)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}
}
