/*
Package filters provides functionality for managing and validating filter parameters for data retrieval.
It includes structures for specifying filters, functions for validating and calculating pagination metadata,
and methods for handling sorting and pagination details.
*/

package filters

import (
	"github.com/k4sper1love/watchlist-api/pkg/validator"
	"math"
	"strings"
)

// Filters represents the parameters for filtering and pagination.
// It includes fields for pagination (Page, PageSize), sorting (Sort), and a list of safe sort options (SortSafeList).
type Filters struct {
	Page         int      `json:"page" validate:"gte=1,lte=10000000"` // Current page number, must be between 1 and 10,000,000.
	PageSize     int      `json:"page_size" validate:"gte=1,lte=100"` // Number of items per page, must be between 1 and 100.
	Sort         string   `json:"sort"`                               // Field to sort by, optionally prefixed with '-' for descending order.
	SortSafeList []string `json:"sort_safe_list"`                     // List of allowed sort fields.
}

// Metadata holds pagination information for the result set.
// It includes the current page, page size, first and last page numbers, and total record count.
type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty" example:"2"`   // The current page number.
	PageSize     int `json:"page_size,omitempty" example:"5"`      // The number of items per page.
	FirstPage    int `json:"first_page,omitempty" example:"1"`     // The first page number, usually 1.
	LastPage     int `json:"last_page,omitempty" example:"3"`      // The last page number based on total records and page size.
	TotalRecords int `json:"total_records,omitempty" example:"15"` // The total number of records available.
}

// CalculateMetadata calculates pagination metadata based on total records, current page, and page size.
func CalculateMetadata(totalRecords, page, pageSize int) Metadata {
	if totalRecords == 0 {
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

// ValidateFilters validates the filter parameters against predefined constraints.
func ValidateFilters(f Filters) (map[string]string, error) {
	errs := validator.ValidateStruct(&f)

	if !isValueInList(f.Sort, f.SortSafeList) {
		if errs == nil {
			errs = make(map[string]string)
		}

		errs["sort"] = "invalid sort value"
	}

	return errs, nil
}

// SortColumn returns the column to sort by, stripping any leading '-' for descending order.
func (f Filters) SortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

// SortDirection returns the sort direction based on whether the Sort field starts with '-'.
// Returns "DESC" for descending order, and "ASC" for ascending order.
func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

// Limit returns the maximum number of items per page (page size).
func (f Filters) Limit() int {
	return f.PageSize
}

// Offset calculates the starting index for the current page of results.
// It is based on the page number and page size.
func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

// isValueInList checks if a given string value is present in a list of strings.
// Returns true if the value is found, otherwise false.
func isValueInList(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}
