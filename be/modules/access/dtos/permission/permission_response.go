package permission

import "harun1804/e-commerce/modules/access/models"

type PermissionResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewPermissionResponseList(permission models.Permission) PermissionResponse {
	return PermissionResponse{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
}

func NewPermissionResponse(permission *models.Permission) PermissionResponse {
	return PermissionResponse{
		Id:          permission.ID,
		Name:        permission.Name,
		Description: permission.Description,
	}
}
