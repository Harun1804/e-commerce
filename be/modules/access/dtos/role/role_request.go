package role

import "harun1804/e-commerce/pkg/query"

type RoleFilterSearch struct {
	query.BaseFilter
}

type RoleSearchRequest struct {
	Search    string `json:"search" validate:"omitempty"`
	Page      int    `json:"page" validate:"omitempty,min=1"`
	Limit     int    `json:"limit" validate:"omitempty"`
	SortBy    string `json:"sortBy" validate:"omitempty,oneof=id name description created_at updated_at"`
	SortOrder string `json:"sortOrder" validate:"omitempty,oneof=asc desc"`
}

type RoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"omitempty,max=255"`
}

var roleAllowedSort = map[string]bool{
	"id":          true,
	"name":        true,
	"description": true,
	"created_at":  true,
	"updated_at":  true,
}

func NewRoleFilter(req RoleSearchRequest) RoleFilterSearch {
	filter := RoleFilterSearch{
		BaseFilter: query.BaseFilter{
			Page:      req.Page,
			Limit:     req.Limit,
			SortBy:    req.SortBy,
			SortOrder: req.SortOrder,
			Search:    req.Search,
		},
	}

	filter.Normalize("id", roleAllowedSort)
	return filter
}
