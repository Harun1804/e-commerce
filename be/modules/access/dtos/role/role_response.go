package role

import (
	permissionDto "harun1804/e-commerce/modules/access/dtos/permission"
	"harun1804/e-commerce/modules/access/models"
)

type RoleResponse struct {
	Id          uint                               `json:"id"`
	Name        string                             `json:"name"`
	Description string                             `json:"description"`
	Permissions []permissionDto.PermissionResponse `json:"permissions,omitempty"`
}

func NewRoleResponseList(role models.Role) RoleResponse {
	return RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

func NewRoleResponse(role *models.Role) RoleResponse {
	resp := RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}

	if len(role.Permissions) > 0 {
		for _, permission := range role.Permissions {
			resp.Permissions = append(resp.Permissions, permissionDto.NewPermissionResponseList(permission))
		}
	}

	return resp
}
