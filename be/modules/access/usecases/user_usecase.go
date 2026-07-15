package usecases

import (
	"context"
	"errors"
	"harun1804/e-commerce/modules/access/dtos/user"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/repositories"
	"harun1804/e-commerce/pkg/hash"
	"harun1804/e-commerce/pkg/httpresponse"
	"harun1804/e-commerce/pkg/logger"

	"gorm.io/gorm"
)

type UserUsecaseInterface interface {
	GetAllUsers(ctx context.Context, filter user.UserFilterSearch) ([]models.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User) error
	UpdateUser(ctx context.Context, user models.User) error
	DeleteUser(ctx context.Context, id uint) error
}

type UserUsecase struct {
	userRepo repositories.UserRepositoryInterface
	entity   string
}

func NewUserUsecase(userRepo repositories.UserRepositoryInterface) UserUsecaseInterface {
	return &UserUsecase{
		userRepo: userRepo,
		entity:   "User",
	}
}

// GetAllUsers implements [UserUsecaseInterface].
func (u *UserUsecase) GetAllUsers(ctx context.Context, filter user.UserFilterSearch) ([]models.User, int64, error) {
	return u.userRepo.GetAllUsers(ctx, filter)
}

// GetUserByID implements [UserUsecaseInterface].
func (u *UserUsecase) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	user, err := u.userRepo.GetUserByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.FailIfError(1, err)
		return nil, errors.New(httpresponse.NotFoundMessage(u.entity, "id", id))
	}
	return user, err
}

// GetUserByUsername implements [UserUsecaseInterface].
func (u *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	user, err := u.userRepo.GetUserByUsername(ctx, username)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.FailIfError(1, err)
		return nil, errors.New(httpresponse.NotFoundMessage(u.entity, "username", username))
	}
	return user, err
}

// CreateUser implements [UserUsecaseInterface].
func (u *UserUsecase) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := hash.GenerateHashPassword(user.Password)
	if err != nil {
		logger.FailIfError(1, err)
		return err
	}
	user.Password = hashedPassword

	return u.userRepo.CreateUser(ctx, user)
}

// UpdateUser implements [UserUsecaseInterface].
func (u *UserUsecase) UpdateUser(ctx context.Context, user models.User) error {
	if user.Password != "" {
		hashedPassword, err := hash.GenerateHashPassword(user.Password)
		if err != nil {
			logger.FailIfError(1, err)
			return err
		}
		user.Password = hashedPassword
	}

	err := u.userRepo.UpdateUser(ctx, user)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.FailIfError(2, err)
		return errors.New(httpresponse.NotFoundMessage(u.entity, "id", user.ID))
	}

	return err
}

// DeleteUser implements [UserUsecaseInterface].
func (u *UserUsecase) DeleteUser(ctx context.Context, id uint) error {
	err := u.userRepo.DeleteUser(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.FailIfError(1, err)
		return errors.New(httpresponse.NotFoundMessage(u.entity, "id", id))
	}

	return err
}
