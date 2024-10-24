package data

import (
	"testing"

	"github.com/2zqa/ssot-specs-server/internal/validator"
)

func TestValidateFilters(t *testing.T) {
	cases := []struct {
		name     string
		filters  Filters
		expected bool
	}{
		{
			name: "valid filters",
			filters: Filters{
				Page:     1,
				PageSize: 10,
				Sort:     "name",
				SortSafelist: []string{
					"name",
					"-name",
				},
			},
			expected: true,
		},
		{
			name: "too low page",
			filters: Filters{
				Page:     0,
				PageSize: 10,
				Sort:     "name",
				SortSafelist: []string{
					"name",
					"-name",
				},
			},
			expected: false,
		},
		{
			name: "too high page",
			filters: Filters{
				Page:     10_000_001,
				PageSize: 10,
				Sort:     "name",
				SortSafelist: []string{
					"name",
					"-name",
				},
			},
			expected: false,
		},
		{
			name: "too low pagesize",
			filters: Filters{
				Page:     1,
				PageSize: 0,
				Sort:     "name",
				SortSafelist: []string{
					"name",
					"-name",
				},
			},
			expected: false,
		},
		{
			name: "too high pagesize",
			filters: Filters{
				Page:     1,
				PageSize: 101,
				Sort:     "name",
				SortSafelist: []string{
					"name",
					"-name",
				},
			},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v := validator.New()
			ValidateFilters(v, tc.filters)
			if v.Valid() != tc.expected {
				t.Errorf("expected %v, but got %v", tc.expected, v.Valid())
			}
		})
	}

}
