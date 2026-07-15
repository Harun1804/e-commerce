package controllers

import (
	"errors"
	"harun1804/e-commerce/modules/access/dtos/user"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/usecases"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/httpresponse"
	"harun1804/e-commerce/pkg/logger"
	"harun1804/e-commerce/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type UserControllerInterface interface {
	GetAllUsers(c fiber.Ctx) error
	GetUserByID(c fiber.Ctx) error
	CreateUser(c fiber.Ctx) error
	UpdateUser(c fiber.Ctx) error
	DeleteUser(c fiber.Ctx) error
}

type UserController struct {
	userUsecase usecases.UserUsecaseInterface
	entity      string
}

func NewUserController(userUsecase usecases.UserUsecaseInterface) UserControllerInterface {
	return &UserController{
		userUsecase: userUsecase,
		entity:      "User",
	}
}

// GetAllUsers implements [UserControllerInterface].
func (u *UserController) GetAllUsers(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	var userSearchReq user.UserSearchRequest
	if err := c.Bind().Query(&userSearchReq); err != nil {
		logger.FailIfError(1, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.BadRequest(c, details)
	}

	filter := user.NewUserFilter(userSearchReq)
	users, totalData, err := u.userUsecase.GetAllUsers(ctx, filter)
	if err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	var userResponse []user.UserResponse
	for _, userModel := range users {
		userResp := user.NewUserResponseList(userModel)
		userResponse = append(userResponse, userResp)
	}

	return httpresponse.SuccessWithPagination(c, httpresponse.FetchMessage(u.entity, true), userResponse, userSearchReq.Page, userSearchReq.Limit, totalData)
}

// GetUserByID implements [UserControllerInterface].
func (u *UserController) GetUserByID(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Params("id")
	if userID == "" {
		logger.FailIfError(1, errors.New("id Parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id Parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(userID)
	userModel, err := u.userUsecase.GetUserByID(ctx, id)
	if err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	userResp := user.NewUserResponse(userModel)
	return httpresponse.Success(c, httpresponse.FetchMessage(u.entity, true), userResp)
}

// CreateUser implements [UserControllerInterface].
func (u *UserController) CreateUser(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	req := user.UserRequest{}
	if err := c.Bind().Body(&req); err != nil {
		logger.FailIfError(1, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.BadRequest(c, details)
	}

	if err := validator.Validate(&req); err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.ValidationError(c, details)
	}

	reqModel := models.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := u.userUsecase.CreateUser(ctx, reqModel); err != nil {
		logger.FailIfError(3, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.CreateMessage(u.entity, true), nil)
}

// UpdateUser implements [UserControllerInterface].
func (u *UserController) UpdateUser(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Params("id")
	if userID == "" {
		logger.FailIfError(1, errors.New("id Parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id Parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(userID)
	req := user.UserRequest{}
	if err := c.Bind().Body(&req); err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.BadRequest(c, details)
	}

	if err := validator.Validate(&req); err != nil {
		logger.FailIfError(3, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.ValidationError(c, details)
	}

	reqModel := models.User{
		ID:       id,
		Username: req.Username,
	}

	if req.Password != "" {
		reqModel.Password = req.Password
	}

	if err := u.userUsecase.UpdateUser(ctx, reqModel); err != nil {
		logger.FailIfError(4, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.UpdateMessage(u.entity, true), nil)
}

// DeleteUser implements [UserControllerInterface].
func (u *UserController) DeleteUser(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	userID := c.Params("id")

	if userID == "" {
		logger.FailIfError(1, errors.New("id Parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id Parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(userID)

	if err := u.userUsecase.DeleteUser(ctx, id); err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.DeleteMessage(u.entity, true), nil)
}
