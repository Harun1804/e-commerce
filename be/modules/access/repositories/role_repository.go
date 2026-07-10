package repositories

import (
	"context"
	"harun1804/e-commerce/modules/access/dtos/role"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error)
	GetRoleByID(ctx context.Context, id uint) (*models.Role, error)
	CreateRole(ctx context.Context, role models.Role) error
	UpdateRole(ctx context.Context, role models.Role) error
	DeleteRole(ctx context.Context, id uint) error
}

type RoleRepository struct {
	db *gorm.DB
}

var roleAttrFields = []string{"id", "name", "description"}

func NewRoleRepository(db *gorm.DB) RoleRepositoryInterface {
	return &RoleRepository{db: db}
}

// GetAllRoles implements [RoleRepositoryInterface].
func (r *RoleRepository) GetAllRoles(ctx context.Context, filter role.RoleFilterSearch) ([]models.Role, int64, error) {
	query := r.db.WithContext(ctx).Model(&models.Role{})

	if filter.HasSearch() {
		search := "%" + filter.Search + "%"
		query = query.Where("name ILIKE ?", search)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, 0, err
	}

	modelRoles := []models.Role{}
	if err := query.
		Select(roleAttrFields).
		Order(filter.Order()).
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&modelRoles).Error; err != nil {
		logger.FailIfError(2, err)
		return nil, 0, err
	}

	return modelRoles, totalRecords, nil
}

// GetRoleByID implements [RoleRepositoryInterface].
func (r *RoleRepository) GetRoleByID(ctx context.Context, id uint) (*models.Role, error) {
	roleModel := &models.Role{}
	if err := r.db.WithContext(ctx).
		Select(roleAttrFields).
		First(roleModel, id).
		Error; err != nil {
		logger.FailIfError(1, err)
		return nil, err
	}

	return roleModel, nil
}

// CreateRole implements [RoleRepositoryInterface].
func (r *RoleRepository) CreateRole(ctx context.Context, role models.Role) error {
	if err := r.db.WithContext(ctx).Create(&role).Error; err != nil {
		logger.FailIfError(1, err)
		return err
	}

	return nil
}

// UpdateRole implements [RoleRepositoryInterface].
func (r *RoleRepository) UpdateRole(ctx context.Context, role models.Role) error {
	result := r.db.WithContext(ctx).Model(&models.Role{}).Where("id = ?", role.ID).Updates(models.Role{
		Name:        role.Name,
		Description: role.Description,
	})

	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", role.ID)); err != nil {
		return err
	}

	return nil
}

// DeleteRole implements [RoleRepositoryInterface].
func (r *RoleRepository) DeleteRole(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Role{}, id)
	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", id)); err != nil {
		return err
	}

	return nil
}
