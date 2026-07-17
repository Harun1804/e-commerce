package usecases

import (
	"context"
	"errors"
	"harun1804/e-commerce/modules/access/dtos/user"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/repositories"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/hash"
	"harun1804/e-commerce/pkg/httpresponse"
	"harun1804/e-commerce/pkg/logger"

	"gorm.io/gorm"
)

type UserUsecaseInterface interface {
	GetAllUsers(ctx context.Context, filter user.UserFilterSearch) ([]models.User, int64, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, user models.User, roleIDs []uint) error
	UpdateUser(ctx context.Context, user models.User, roleIDs []uint) error
	DeleteUser(ctx context.Context, id uint) error
	AttachRoleToUser(ctx context.Context, userID uint, roleIDs []uint) error
	DetachRoleFromUser(ctx context.Context, userID uint, roleIDs []uint) error
	SyncUserRoles(ctx context.Context, userID uint, roleIDs []uint) error
}

type UserUsecase struct {
	userRepo     repositories.UserRepositoryInterface
	roleUsecase  RoleUsecaseInterface
	roleUserRepo repositories.RoleUserRepositoryInterface
	entity       string
}

func NewUserUsecase(
	userRepo repositories.UserRepositoryInterface,
	roleUsecase RoleUsecaseInterface,
	roleUserRepo repositories.RoleUserRepositoryInterface,
) UserUsecaseInterface {
	return &UserUsecase{
		userRepo:     userRepo,
		roleUsecase:  roleUsecase,
		roleUserRepo: roleUserRepo,
		entity:       "User",
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
func (u *UserUsecase) CreateUser(ctx context.Context, user models.User, roleIDs []uint) error {
	hashedPassword, err := hash.GenerateHashPassword(user.Password)
	if err != nil {
		logger.FailIfError(1, err)
		return err
	}
	user.Password = hashedPassword

	if err := u.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	if len(roleIDs) > 0 {
		user, err := u.userRepo.GetUserByUsername(ctx, user.Username)
		if err != nil {
			logger.FailIfError(2, err)
			return err
		}
		return u.AttachRoleToUser(ctx, user.ID, roleIDs)
	}

	return nil
}

// UpdateUser implements [UserUsecaseInterface].
func (u *UserUsecase) UpdateUser(ctx context.Context, user models.User, roleIDs []uint) error {
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

	if len(roleIDs) > 0 {
		if err := u.SyncUserRoles(ctx, user.ID, roleIDs); err != nil {
			return err
		}
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

// AttachRoleToUser implements [UserUsecaseInterface].
func (u *UserUsecase) AttachRoleToUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := u.validateUserRoles(ctx, userID, roleIDs); err != nil {
		return err
	}

	return u.roleUserRepo.AttachRoleUser(ctx, userID, conv.UniqueValues(roleIDs))
}

// DetachRoleFromUser implements [UserUsecaseInterface].
func (u *UserUsecase) DetachRoleFromUser(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := u.validateUserRoles(ctx, userID, roleIDs); err != nil {
		return err
	}

	return u.roleUserRepo.DetachRoleUser(ctx, userID, conv.UniqueValues(roleIDs))
}

// SyncUserRoles makes a user's roles exactly match the provided role IDs.
func (u *UserUsecase) SyncUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if err := u.validateUserRoles(ctx, userID, roleIDs); err != nil {
		return err
	}

	uniqueRoleIDs := conv.UniqueValues(roleIDs)
	existingRoleIDs, err := u.roleUserRepo.GetRoleIDsByUserID(ctx, userID)
	if err != nil {
		return err
	}

	roleIDsToAttach := conv.MissingValues(uniqueRoleIDs, existingRoleIDs)
	roleIDsToDetach := conv.MissingValues(existingRoleIDs, uniqueRoleIDs)

	if err := u.roleUserRepo.AttachRoleUser(ctx, userID, roleIDsToAttach); err != nil {
		return err
	}

	return u.roleUserRepo.DetachRoleUser(ctx, userID, roleIDsToDetach)
}

func (u *UserUsecase) validateUserRoles(ctx context.Context, userID uint, roleIDs []uint) error {
	if _, err := u.GetUserByID(ctx, userID); err != nil {
		return err
	}

	uniqueRoleIDs := conv.UniqueValues(roleIDs)
	roles, err := u.roleUsecase.GetRolesByIDs(ctx, uniqueRoleIDs, false)
	if err != nil {
		return err
	}

	if len(roles) != len(uniqueRoleIDs) {
		return errors.New(httpresponse.NotFoundMessage("Role", "id", conv.MissingValues(uniqueRoleIDs, extractRoleIDs(roles))))
	}

	return nil
}

func extractRoleIDs(roles []models.Role) []uint {
	ids := make([]uint, 0, len(roles))
	for _, role := range roles {
		ids = append(ids, role.ID)
	}
	return ids
}
