package repositories

import (
	"context"
	"harun1804/e-commerce/modules/access/models/pivot"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RolePermissionRepositoryInterface interface {
	GetPermissionIDsByRoleID(ctx context.Context, roleId uint) ([]uint, error)
	AttachRolePermission(ctx context.Context, roleId uint, permissionIds []uint) error
	DetachRolePermission(ctx context.Context, roleId uint, permissionIds []uint) error
}

type RolePermissionRepository struct {
	db *gorm.DB
}

func NewRolePermissionRepository(db *gorm.DB) RolePermissionRepositoryInterface {
	return &RolePermissionRepository{
		db: db,
	}
}

// GetPermissionIDsByRoleID implements [RolePermissionRepositoryInterface].
func (r *RolePermissionRepository) GetPermissionIDsByRoleID(ctx context.Context, roleId uint) ([]uint, error) {
	permissionIds := []uint{}
	if err := r.db.WithContext(ctx).
		Model(&pivot.RolePermission{}).
		Where("role_id = ?", roleId).
		Pluck("permission_id", &permissionIds).Error; err != nil {
		return nil, err
	}

	return permissionIds, nil
}

// AttachRolePermission implements [RolePermissionRepositoryInterface].
func (r *RolePermissionRepository) AttachRolePermission(ctx context.Context, roleId uint, permissionIds []uint) error {
	var rolePermissions []pivot.RolePermission
	for _, permissionId := range permissionIds {
		rolePermissions = append(rolePermissions, pivot.RolePermission{
			RoleId:       roleId,
			PermissionId: permissionId,
		})
	}

	if len(rolePermissions) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&rolePermissions).Error
}

// DetachRolePermission implements [RolePermissionRepositoryInterface].
func (r *RolePermissionRepository) DetachRolePermission(ctx context.Context, roleId uint, permissionIds []uint) error {
	if len(permissionIds) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Where("role_id = ? AND permission_id IN ?", roleId, permissionIds).Delete(&pivot.RolePermission{}).Error
}
