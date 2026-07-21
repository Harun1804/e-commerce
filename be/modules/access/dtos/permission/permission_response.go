package permission

import "harun1804/e-commerce/modules/access/models"

type PermissionResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Module      string `json:"module"`
	Description string `json:"description"`
}

func NewPermissionResponseList(permission models.Permission) PermissionResponse {
	return PermissionResponse{
		Id:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Module:      permission.Module,
		Description: permission.Description,
	}
}

func NewPermissionResponse(permission *models.Permission) PermissionResponse {
	return PermissionResponse{
		Id:          permission.ID,
		Name:        permission.Name,
		Slug:        permission.Slug,
		Module:      permission.Module,
		Description: permission.Description,
	}
}
