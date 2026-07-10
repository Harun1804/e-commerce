package query

import "strings"

type BaseFilter struct {
	Page      int
	Limit     int
	SortBy    string
	SortOrder string
	Search    string
}

func (f *BaseFilter) Normalize(defaultSort string, allowedSort map[string]bool) {
	if f.Page <= 0 {
		f.Page = 1
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}

	if f.SortBy == "" {
		f.SortBy = defaultSort
	}

	// whitelist sort field
	if !allowedSort[f.SortBy] {
		f.SortBy = defaultSort
	}

	if f.SortOrder == "" {
		f.SortOrder = "desc"
	}

	f.SortOrder = strings.ToLower(f.SortOrder)
	if f.SortOrder != "asc" && f.SortOrder != "desc" {
		f.SortOrder = "desc"
	}

	f.Search = strings.TrimSpace(f.Search)
}

func (f *BaseFilter) Offset() int {
	return (f.Page - 1) * f.Limit
}

func (f *BaseFilter) Order() string {
	return f.SortBy + " " + f.SortOrder
}

func (f *BaseFilter) HasSearch() bool {
	return f.Search != ""
}