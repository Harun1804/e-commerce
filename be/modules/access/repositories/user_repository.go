package repositories

import (
	"context"
	"harun1804/e-commerce/modules/access/dtos/user"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetAllUsers(ctx context.Context, filter user.UserFilterSearch) ([]models.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserRepository struct {
	db *gorm.DB
}

var userAttrFields = []string{"id", "username", "created_at"}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

// GetAllUsers implements [UserRepositoryInterface].
func (u *UserRepository) GetAllUsers(ctx context.Context, filter user.UserFilterSearch) ([]models.User, int64, error) {
	query := u.db.WithContext(ctx).Model(&models.User{})

	if filter.HasSearch() {
		search := "%" + filter.Search + "%"
		query = query.Where("username ILIKE ?", search)
	}

	var totalRecords int64
	if err := query.Count(&totalRecords).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, 0, err
	}

	modelUsers := []models.User{}
	if err := query.
		Select(userAttrFields).
		Order(filter.Order()).
		Limit(filter.Limit).
		Offset(filter.Offset()).
		Find(&modelUsers).Error; err != nil {
		logger.FailIfError(2, err)
		return nil, 0, err
	}

	return modelUsers, totalRecords, nil
}

// GetUserByID implements [UserRepositoryInterface].
func (u *UserRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	userModel := &models.User{}
	if err := u.db.WithContext(ctx).First(userModel, id).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, err
	}

	return userModel, nil
}

// GetUserByUsername implements [UserRepositoryInterface].
func (u *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	userModel := &models.User{}
	if err := u.db.WithContext(ctx).Where("username = ?", username).First(userModel).Error; err != nil {
		logger.FailIfError(1, err)
		return nil, err
	}

	return userModel, nil
}

// CreateUser implements [UserRepositoryInterface].
func (u *UserRepository) CreateUser(ctx context.Context, user models.User) error {
	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		logger.FailIfError(1, err)
		return err
	}

	return nil
}

// UpdateUser implements [UserRepositoryInterface].
func (u *UserRepository) UpdateUser(ctx context.Context, user models.User) error {
	result := u.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", user.ID).Updates(user)
	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", user.ID)); err != nil {
		return err
	}

	return nil
}

// DeleteUser implements [UserRepositoryInterface].
func (u *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	result := u.db.WithContext(ctx).Delete(&models.User{}, id)
	if err := logger.FailIfError(1, result.Error); err != nil {
		return err
	}

	if err := logger.FailIfRowsAffectedZero(2, result.RowsAffected, zap.Uint("id", id)); err != nil {
		return err
	}

	return nil
}
