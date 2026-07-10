package usecases

import (
	"context"
	"errors"
	"harun1804/e-commerce/modules/access/dtos/role"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/repositories"
	"harun1804/e-commerce/pkg/httpresponse"

	"gorm.io/gorm"
)

type RoleUsecaseInterface interface {
	GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error)
	GetRoleByID(ctx context.Context, roleID uint) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, roleID uint) error
}

type roleUsecase struct {
	roleRepo repositories.RoleRepositoryInterface
	entity   string
}

func NewRoleUsecase(roleRepo repositories.RoleRepositoryInterface) RoleUsecaseInterface {
	return &roleUsecase{
		roleRepo: roleRepo,
		entity:   "Role",
	}
}

// GetAllRoles implements RoleUsecaseInterface.
func (r *roleUsecase) GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error) {
	return r.roleRepo.GetAllRoles(ctx, filter)
}

// GetRoleByID implements RoleUsecaseInterface.
func (r *roleUsecase) GetRoleByID(ctx context.Context, roleID uint) (*models.Role, error) {
	return r.roleRepo.GetRoleByID(ctx, roleID)
}

// CreateRole implements RoleUsecaseInterface.
func (r *roleUsecase) CreateRole(ctx context.Context, role models.Role) error {
	return r.roleRepo.CreateRole(ctx, role)
}

// UpdateRole implements RoleUsecaseInterface.
func (r *roleUsecase) UpdateRole(ctx context.Context, role models.Role) error {
	err := r.roleRepo.UpdateRole(ctx, role)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(httpresponse.NotFoundMessage(r.entity, "id", role.ID))
	}
	return err
}

// DeleteRole implements RoleUsecaseInterface.
func (r *roleUsecase) DeleteRole(ctx context.Context, roleID uint) error {
	err := r.roleRepo.DeleteRole(ctx, roleID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(httpresponse.NotFoundMessage(r.entity, "id", roleID))
	}
	return err
}
