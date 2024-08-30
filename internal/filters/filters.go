package filters

import (
	"github.com/k4sper1love/watchlist-api/internal/validator"
	"math"
	"strings"
)

type Filters struct {
	Page         int      `json:"page" validate:"gte=1,lte=10000000"`
	PageSize     int      `json:"page_size" validate:"gte=1,lte=100"`
	Sort         string   `json:"sort"`
	SortSafeList []string `json:"sort_safe_list"`
}

type Metadata struct {
	CurrentPage  int `json:"current_page,omitempty"`
	PageSize     int `json:"page_size,omitempty"`
	FirstPage    int `json:"first_page,omitempty"`
	LastPage     int `json:"last_page,omitempty"`
	TotalRecords int `json:"total_records,omitempty"`
}

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

func ValidateFilters(f Filters) map[string]string {
	errs := validator.ValidateStruct(&f)

	if !isValueInList(f.Sort, f.SortSafeList) {
		if errs == nil {
			errs = make(map[string]string)
		}

		errs["sort"] = "invalid sort value"
	}

	return errs
}

func (f Filters) SortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

func (f Filters) SortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}

func (f Filters) Limit() int {
	return f.PageSize
}

func (f Filters) Offset() int {
	return (f.Page - 1) * f.PageSize
}

func isValueInList(s string, list []string) bool {
	for _, v := range list {
		if s == v {
			return true
		}
	}
	return false
}
