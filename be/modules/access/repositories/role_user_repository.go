package repositories

import (
	"context"
	"harun1804/e-commerce/modules/access/models/pivot"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RoleUserRepositoryInterface interface {
	GetRoleIDsByUserID(ctx context.Context, userId uint) ([]uint, error)
	AttachRoleUser(ctx context.Context, userId uint, roleIds []uint) error
	DetachRoleUser(ctx context.Context, userId uint, roleIds []uint) error
}

type RoleUserRepository struct {
	db *gorm.DB
}

func NewRoleUserRepository(db *gorm.DB) RoleUserRepositoryInterface {
	return &RoleUserRepository{
		db: db,
	}
}

// GetRoleIDsByUserID implements [RoleUserRepositoryInterface].
func (r *RoleUserRepository) GetRoleIDsByUserID(ctx context.Context, userId uint) ([]uint, error) {
	roleIds := []uint{}
	if err := r.db.WithContext(ctx).
		Model(&pivot.RoleUser{}).
		Where("user_id = ?", userId).
		Pluck("role_id", &roleIds).Error; err != nil {
		return nil, err
	}

	return roleIds, nil
}

// AttachRoleUser implements [RoleUserRepositoryInterface].
func (r *RoleUserRepository) AttachRoleUser(ctx context.Context, userId uint, roleIds []uint) error {
	var roleUsers []pivot.RoleUser
	for _, roleId := range roleIds {
		roleUsers = append(roleUsers, pivot.RoleUser{
			RoleId: roleId,
			UserId: userId,
		})
	}

	if len(roleUsers) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Clauses(clause.OnConflict{DoNothing: true}).Create(&roleUsers).Error
}

// DetachRoleUser implements [RoleUserRepositoryInterface].
func (r *RoleUserRepository) DetachRoleUser(ctx context.Context, userId uint, roleIds []uint) error {
	if len(roleIds) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Where("user_id = ? AND role_id IN ?", userId, roleIds).Delete(&pivot.RoleUser{}).Error
}
