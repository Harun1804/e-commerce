package usecases

import (
	"context"
	"errors"
	"harun1804/e-commerce/modules/access/dtos/permission"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/repositories"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/httpresponse"

	"gorm.io/gorm"
)

type PermissionUsecaseInterface interface {
	GetAllPermissions(ctx context.Context, filter permission.PermissionFilterSearch) ([]models.Permission, int64, error)
	GetPermissionByID(ctx context.Context, id uint) (*models.Permission, error)
	GetPermissionsByIDs(ctx context.Context, ids []uint) ([]models.Permission, error)
	CreatePermission(ctx context.Context, permission models.Permission) error
	UpdatePermission(ctx context.Context, permission models.Permission) error
	DeletePermission(ctx context.Context, id uint) error
}

type permissionUsecase struct {
	permissionRepo repositories.PermissionRepositoryInterface
	entity         string
}

func NewPermissionUsecase(permissionRepo repositories.PermissionRepositoryInterface) PermissionUsecaseInterface {
	return &permissionUsecase{
		permissionRepo: permissionRepo,
		entity:         "Permission",
	}
}

// GetAllPermissions implements PermissionUsecaseInterface.
func (p *permissionUsecase) GetAllPermissions(ctx context.Context, filter permission.PermissionFilterSearch) ([]models.Permission, int64, error) {
	return p.permissionRepo.GetAllPermissions(ctx, filter)
}

// GetPermissionByID implements PermissionUsecaseInterface.
func (p *permissionUsecase) GetPermissionByID(ctx context.Context, id uint) (*models.Permission, error) {
	return p.permissionRepo.GetPermissionByID(ctx, id)
}

// GetPermissionsByIDs implements [PermissionUsecaseInterface].
func (p *permissionUsecase) GetPermissionsByIDs(ctx context.Context, ids []uint) ([]models.Permission, error) {
	return p.permissionRepo.GetPermissionsByIDs(ctx, ids)
}

// CreatePermission implements PermissionUsecaseInterface.
func (p *permissionUsecase) CreatePermission(ctx context.Context, permission models.Permission) error {
	permission.Slug = conv.GenerateSlug(permission.Name, true)
	return p.permissionRepo.CreatePermission(ctx, permission)
}

// UpdatePermission implements PermissionUsecaseInterface.
func (p *permissionUsecase) UpdatePermission(ctx context.Context, permission models.Permission) error {
	if permission.Name != "" {
		permission.Slug = conv.GenerateSlug(permission.Name, true)
	}

	err := p.permissionRepo.UpdatePermission(ctx, permission)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(httpresponse.NotFoundMessage(p.entity, "id", permission.ID))
	}
	return err
}

// DeletePermission implements PermissionUsecaseInterface.
func (p *permissionUsecase) DeletePermission(ctx context.Context, id uint) error {
	err := p.permissionRepo.DeletePermission(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New(httpresponse.NotFoundMessage(p.entity, "id", id))
	}
	return err
}
