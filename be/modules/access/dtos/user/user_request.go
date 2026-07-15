package user

import "harun1804/e-commerce/pkg/query"

type UserFilterSearch struct {
	query.BaseFilter
}

type UserSearchRequest struct {
	Search    string `json:"search" validate:"omitempty"`
	Page      int    `json:"page" validate:"omitempty,min=1"`
	Limit     int    `json:"limit" validate:"omitempty"`
	SortBy    string `json:"sortBy" validate:"omitempty,oneof=id username created_at updated_at"`
	SortOrder string `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

type UserRequest struct {
	Username string `json:"username" binding:"required" validate:"required,min=3,max=100"`
	Password string `json:"password" binding:"required" validate:"omitempty,min=5"`
}

var userAllowedSort = map[string]bool{
	"id":          true,
	"username":    true,
	"created_at":  true,
	"updated_at":  true,
}

func NewUserFilter(req UserSearchRequest) UserFilterSearch {
	filter := UserFilterSearch{
		BaseFilter: query.BaseFilter{
			Page:      req.Page,
			Limit:     req.Limit,
			SortBy:    req.SortBy,
			SortOrder: req.SortOrder,
			Search:    req.Search,
		},
	}

	filter.Normalize("id", userAllowedSort)
	return filter
}