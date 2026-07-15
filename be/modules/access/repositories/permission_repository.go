package repositories

import (
	"context"
	"harun1804/e-commerce/modules/access/dtos/permission"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PermissionRepositoryInterface interface {
	GetAllPermissions(ctx context.Context, filter permission.PermissionFilterSearch) ([]models.Permission, int64, error)
	GetPermissionByID(ctx context.Context, id uint) (*models.Permission, error)
	CreatePermission(ctx context.Context, permission models.Permission) error
	UpdatePermission(ctx context.Context, permission models.Permission) error
	DeletePermission(ctx context.Context, id uint) error
}

type PermissionRepository struct {
	db *gorm.DB
}

var permissionAttrFields = []string{"id", "name", "description"}

func NewPermissionRepository(db *gorm.DB) PermissionRepositoryInterface {
	return &PermissionRepository{db: db}
}

// GetAllPermissions implements [PermissionRepositoryInterface].
func (p *PermissionRepository) GetAllPermissions(ctx context.Context, filter permission.PermissionFilterSearch) ([]models.Permission, int64, error) {
	query := p.db.WithContext(ctx).Model(&models.Permission{})

	if filter.HasSearch() {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, 0, err
	}

	modelPermissions := []models.Permission{}
	if err := query.
		Select(permissionAttrFields).
		Order(filter.Order()).
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&modelPermissions).Error; err != nil {
		logger.FailIfError(2, err)
		return nil, 0, err
	}

	return modelPermissions, totalRecords, nil
}

// GetPermissionByID implements [PermissionRepositoryInterface].
func (p *PermissionRepository) GetPermissionByID(ctx context.Context, id uint) (*models.Permission, error) {
	permissionModel := &models.Permission{}
	if err := p.db.WithContext(ctx).
		Select(permissionAttrFields).
		First(permissionModel, id).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, err
	}
	return permissionModel, nil
}

// CreatePermission implements [PermissionRepositoryInterface].
func (p *PermissionRepository) CreatePermission(ctx context.Context, permission models.Permission) error {
	if err := p.db.WithContext(ctx).Create(&permission).Error; err != nil {
		logger.FailIfError(1, err)
		return err
	}
	return nil
}

// UpdatePermission implements [PermissionRepositoryInterface].
func (p *PermissionRepository) UpdatePermission(ctx context.Context, permission models.Permission) error {
	result := p.db.WithContext(ctx).Model(&models.Permission{}).Where("id = ?", permission.ID).Updates(permission)
	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", permission.ID)); err != nil {
		return err
	}

	return nil
}

// DeletePermission implements [PermissionRepositoryInterface].
func (p *PermissionRepository) DeletePermission(ctx context.Context, id uint) error {
	result := p.db.WithContext(ctx).Delete(&models.Permission{}, id)
	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", id)); err != nil {
		return err
	}

	return nil
}
