package data

import (
	"math"
	"strings"

	"github.com/2zqa/ssot-specs-server/internal/validator"
	"gorm.io/gorm"
)

type Filters struct {
	Page         int
	PageSize     int
	Sort         string
	SortSafelist []string
}

// Metadata contains pagination metadata.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

// calculateMetadata calculates the appropriate pagination metadata
// values given the total number of records, current page, and page size values.
func calculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
		// Note that we return an empty Metadata struct if there are no records.
		return Metadata{}
	}

	return Metadata{
		CurrentPage:  page,
		PageSize:     pageSize,
		FirstPage:    1,
		LastPage:     int(math.Ceil(float64(totalRecords) / float64(pageSize))),
		TotalRecords: totalRecords,
	}
}

// sortColumn checks that the client-provided Sort field matches one of the entries in the safelist
// and if it does, extract the column name from the Sort field by stripping the leading
// hyphen character (if one exists).
func (f Filters) sortColumn() string {
	for _, safeValue := range f.SortSafelist {
		if f.Sort == safeValue {
			return strings.TrimPrefix(f.Sort, "-")
		}
	}

	panic("unsafe sort parameter: " + f.Sort)
}

// sortDirection returns the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}

	return "ASC"
}

func (f Filters) sort() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Order(f.sortColumn() + " " + f.sortDirection()).Order("id")
	}
}

func (f Filters) paginate(db *gorm.DB, searchTerms string, metadata *Metadata) func(db *gorm.DB) *gorm.DB {
	var totalRecords int64
	db.Model(&Device{}).
		Scopes(f.search(searchTerms)).
		Count(&totalRecords)
	*metadata = calculateMetadata(int(totalRecords), f.Page, f.PageSize)
	return func(db *gorm.DB) *gorm.DB {
		offset := (f.Page - 1) * f.PageSize
		return db.Offset(offset).Limit(f.PageSize)
	}
}

func (f Filters) search(searchTerms string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if searchTerms == "" {
			return db
		}
		return db.Joins("JOIN searches ON searches.device_id = devices.id").
			Where("searches.data ILIKE ?", "%"+searchTerms+"%")
	}
}

func ValidateFilters(v *validator.Validator, f Filters) {
	// Check that the page and page_size parameters contain sensible values.
	v.Check(f.Page > 0, "page", "must be greater than zero")
	v.Check(f.Page <= 10_000_000, "page", "must be a maximum of 10 million")
	v.Check(f.PageSize > 0, "page_size", "must be greater than zero")
	v.Check(f.PageSize <= 100, "page_size", "must be a maximum of 100")

	// Check that the sort parameter matches a value in the safelist.
	v.Check(validator.PermittedValue(f.Sort, f.SortSafelist...), "sort", "invalid sort value")
}
