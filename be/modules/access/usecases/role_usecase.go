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
	GetRoleByID(ctx context.Context, roleID uint, needDetailPermission bool) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, roleID uint) error
	AttachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error
	DetachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error
}

type roleUsecase struct {
	roleRepo           repositories.RoleRepositoryInterface
	permissionRepo     repositories.PermissionRepositoryInterface
	rolePermissionRepo repositories.RolePermissionRepositoryInterface
	entity             string
}

func NewRoleUsecase(
	roleRepo repositories.RoleRepositoryInterface,
	permissionRepo repositories.PermissionRepositoryInterface,
	rolePermissionRepo repositories.RolePermissionRepositoryInterface,
) RoleUsecaseInterface {
	return &roleUsecase{
		roleRepo:           roleRepo,
		permissionRepo:     permissionRepo,
		rolePermissionRepo: rolePermissionRepo,
		entity:             "Role",
	}
}

// GetAllRoles implements RoleUsecaseInterface.
func (r *roleUsecase) GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error) {
	return r.roleRepo.GetAllRoles(ctx, filter)
}

// GetRoleByID implements RoleUsecaseInterface.
func (r *roleUsecase) GetRoleByID(ctx context.Context, roleID uint, needDetailPermission bool) (*models.Role, error) {
	role, err := r.roleRepo.GetRoleByID(ctx, roleID, needDetailPermission)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New(httpresponse.NotFoundMessage(r.entity, "id", roleID))
	}
	return role, err
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

// AttachRolePermission assigns permissions to a role.
func (r *roleUsecase) AttachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if err := r.validateRolePermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	return r.rolePermissionRepo.AttachRolePermission(ctx, roleID, permissionIDs)
}

// DetachRolePermission removes permissions from a role.
func (r *roleUsecase) DetachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if err := r.validateRolePermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	return r.rolePermissionRepo.DetachRolePermission(ctx, roleID, permissionIDs)
}

func (r *roleUsecase) validateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if _, err := r.GetRoleByID(ctx, roleID, false); err != nil {
		return err
	}

	uniquePermissionIDs := uniqueUint(permissionIDs)
	permissions, err := r.permissionRepo.GetPermissionsByIDs(ctx, uniquePermissionIDs)
	if err != nil {
		return err
	}

	if len(permissions) != len(uniquePermissionIDs) {
		return errors.New(httpresponse.NotFoundMessage("Permission", "id", missingPermissionIDs(uniquePermissionIDs, permissions)))
	}

	return nil
}

func uniqueUint(values []uint) []uint {
	seen := map[uint]bool{}
	uniqueValues := []uint{}

	for _, value := range values {
		if seen[value] {
			continue
		}
		seen[value] = true
		uniqueValues = append(uniqueValues, value)
	}

	return uniqueValues
}

func missingPermissionIDs(permissionIDs []uint, permissions []models.Permission) []uint {
	found := map[uint]bool{}
	for _, permission := range permissions {
		found[permission.ID] = true
	}

	missingIDs := []uint{}
	for _, permissionID := range permissionIDs {
		if !found[permissionID] {
			missingIDs = append(missingIDs, permissionID)
		}
	}

	return missingIDs
}
