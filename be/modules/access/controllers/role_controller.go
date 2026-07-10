package controllers

import (
	"errors"
	"harun1804/e-commerce/modules/access/dtos/role"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/usecases"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/httpresponse"
	"harun1804/e-commerce/pkg/logger"
	"harun1804/e-commerce/pkg/validator"

	"github.com/gofiber/fiber/v3"
	"go.uber.org/zap"
)

type RoleControllerInterface interface {
	GetAllRoles(c fiber.Ctx) error
	GetRoleByID(c fiber.Ctx) error
	CreateRole(c fiber.Ctx) error
	UpdateRole(c fiber.Ctx) error
	DeleteRole(c fiber.Ctx) error
}

type RoleController struct {
	roleUsecase usecases.RoleUsecaseInterface
	entity      string
}

func NewRoleController(roleUsecase usecases.RoleUsecaseInterface) RoleControllerInterface {
	return &RoleController{
		roleUsecase: roleUsecase,
		entity:      "Role",
	}
}

// GetAllRoles implements [RoleControllerInterface].
func (r *RoleController) GetAllRoles(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	var searchReq role.RoleSearchRequest

	if err := c.Bind().Query(&searchReq); err != nil {
		logger.FailIfError(1, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.BadRequest(c, details)
	}

	zap.L().Info("Received search request: ", zap.Any("searchReq", searchReq))
	filter := role.NewRoleFilter(searchReq)

	roles, totalData, err := r.roleUsecase.GetAllRoles(ctx, filter)
	if err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	var roleResponses []role.RoleResponse
	for _, roleModel := range roles {
		roleResp := role.NewRoleResponseList(roleModel)
		roleResponses = append(roleResponses, roleResp)
	}

	return httpresponse.SuccessWithPagination(c, httpresponse.FetchMessage(r.entity, true), roleResponses, searchReq.Page, searchReq.Limit, totalData)
}

// GetRoleByID implements [RoleControllerInterface].
func (r *RoleController) GetRoleByID(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	roleID := c.Params("id")
	if roleID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(roleID)
	roleModel, err := r.roleUsecase.GetRoleByID(ctx, id)
	if err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	roleResp := role.NewRoleResponse(roleModel)
	return httpresponse.Success(c, httpresponse.FetchMessage(r.entity, true), roleResp)
}

// CreateRole implements [RoleControllerInterface].
func (r *RoleController) CreateRole(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	req := role.RoleRequest{}

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

	reqModel := models.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := r.roleUsecase.CreateRole(ctx, reqModel); err != nil {
		logger.FailIfError(3, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.CreateMessage(r.entity, true), nil)
}

// UpdateRole implements [RoleControllerInterface].
func (r *RoleController) UpdateRole(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	roleID := c.Params("id")
	if roleID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(roleID)
	req := role.RoleRequest{}

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

	reqModel := models.Role{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := r.roleUsecase.UpdateRole(ctx, reqModel); err != nil {
		logger.FailIfError(3, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.UpdateMessage(r.entity, true), nil)
}

// DeleteRole implements [RoleControllerInterface].
func (r *RoleController) DeleteRole(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	roleID := c.Params("id")
	if roleID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(roleID)
	if err := r.roleUsecase.DeleteRole(ctx, id); err != nil {
		logger.FailIfError(2, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.InternalServerError(c, details)
	}

	return httpresponse.Success(c, httpresponse.DeleteMessage(r.entity, true), nil)
}
