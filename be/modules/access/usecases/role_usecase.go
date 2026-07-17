package usecases

import (
	"context"
	"errors"
	"harun1804/e-commerce/modules/access/dtos/role"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/repositories"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/httpresponse"

	"gorm.io/gorm"
)

type RoleUsecaseInterface interface {
	GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error)
	GetRoleByID(ctx context.Context, roleID uint, needDetailPermission bool) (*models.Role, error)
	GetRolesByIDs(ctx context.Context, ids []uint, needDetailPermission bool) ([]models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, roleID uint) error
	AttachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error
	DetachRolePermission(ctx context.Context, roleID uint, permissionIDs []uint) error
	SyncRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error
}

type roleUsecase struct {
	roleRepo           repositories.RoleRepositoryInterface
	rolePermissionRepo repositories.RolePermissionRepositoryInterface
	permissionUsecase  PermissionUsecaseInterface
	entity             string
}

func NewRoleUsecase(
	roleRepo repositories.RoleRepositoryInterface,
	permissionUsecase PermissionUsecaseInterface,
	rolePermissionRepo repositories.RolePermissionRepositoryInterface,
) RoleUsecaseInterface {
	return &roleUsecase{
		roleRepo:           roleRepo,
		permissionUsecase:  permissionUsecase,
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

// GetRolesByIDs implements [RoleUsecaseInterface].
func (r *roleUsecase) GetRolesByIDs(ctx context.Context, ids []uint, needDetailPermission bool) ([]models.Role, error) {
	return r.roleRepo.GetRolesByIDs(ctx, ids, needDetailPermission)
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

// SyncRolePermissions makes a role's permissions exactly match the provided permission IDs.
func (r *roleUsecase) SyncRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if err := r.validateRolePermissions(ctx, roleID, permissionIDs); err != nil {
		return err
	}

	uniquePermissionIDs := conv.UniqueValues(permissionIDs)
	existingPermissionIDs, err := r.rolePermissionRepo.GetPermissionIDsByRoleID(ctx, roleID)
	if err != nil {
		return err
	}

	permissionIDsToAttach := conv.MissingValues(uniquePermissionIDs, existingPermissionIDs)
	permissionIDsToDetach := conv.MissingValues(existingPermissionIDs, uniquePermissionIDs)

	if err := r.rolePermissionRepo.AttachRolePermission(ctx, roleID, permissionIDsToAttach); err != nil {
		return err
	}

	return r.rolePermissionRepo.DetachRolePermission(ctx, roleID, permissionIDsToDetach)
}

func (r *roleUsecase) validateRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if _, err := r.GetRoleByID(ctx, roleID, false); err != nil {
		return err
	}

	uniquePermissionIDs := conv.UniqueValues(permissionIDs)
	permissions, err := r.permissionUsecase.GetPermissionsByIDs(ctx, uniquePermissionIDs)
	if err != nil {
		return err
	}

	if len(permissions) != len(uniquePermissionIDs) {
		return errors.New(httpresponse.NotFoundMessage("Permission", "id", conv.MissingValues(uniquePermissionIDs, extractPermissionIDs(permissions))))
	}

	return nil
}

func extractPermissionIDs(permissions []models.Permission) []uint {
	ids := make([]uint, 0, len(permissions))
	for _, permission := range permissions {
		ids = append(ids, permission.ID)
	}
	return ids
}
