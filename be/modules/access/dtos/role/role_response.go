package role

import "harun1804/e-commerce/modules/access/models"

type RoleResponse struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func NewRoleResponseList(role models.Role) RoleResponse {
	return RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}

func NewRoleResponse(role *models.Role) RoleResponse {
	return RoleResponse{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
	}
}
