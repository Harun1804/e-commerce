package controllers

import (
	"errors"
	"harun1804/e-commerce/modules/access/dtos/permission"
	"harun1804/e-commerce/modules/access/models"
	"harun1804/e-commerce/modules/access/usecases"
	"harun1804/e-commerce/pkg/conv"
	"harun1804/e-commerce/pkg/httpresponse"
	"harun1804/e-commerce/pkg/logger"
	"harun1804/e-commerce/pkg/validator"

	"github.com/gofiber/fiber/v3"
)

type PermissionControllerInterface interface {
	GetAllPermissions(c fiber.Ctx) error
	GetPermissionByID(c fiber.Ctx) error
	CreatePermission(c fiber.Ctx) error
	UpdatePermission(c fiber.Ctx) error
	DeletePermission(c fiber.Ctx) error
}

type PermissionController struct {
	permissionUsecase usecases.PermissionUsecaseInterface
	entity            string
}

func NewPermissionController(permissionUsecase usecases.PermissionUsecaseInterface) PermissionControllerInterface {
	return &PermissionController{
		permissionUsecase: permissionUsecase,
		entity:            "Permission",
	}
}

// GetAllPermissions implements [PermissionControllerInterface].
func (p *PermissionController) GetAllPermissions(c fiber.Ctx) error {
	ctx := c.RequestCtx()
	var searchReq permission.PermissionSearchRequest

	if err := c.Bind().Query(&searchReq); err != nil {
		logger.FailIfError(1, err)
		details := httpresponse.ErrorDetail(err)
		return httpresponse.BadRequest(c, details)
	}

	filter := permission.NewPermissionFilter(searchReq)

	permissions, totalData, err := p.permissionUsecase.GetAllPermissions(ctx, filter)
	if err != nil {
		logger.FailIfError(2, err)
		return httpresponse.Error(c, err)
	}

	var permissionResponses []permission.PermissionResponse
	for _, permissionModel := range permissions {
		permissionResp := permission.NewPermissionResponseList(permissionModel)
		permissionResponses = append(permissionResponses, permissionResp)
	}

	return httpresponse.SuccessWithPagination(c, httpresponse.FetchMessage(p.entity, true), permissionResponses, searchReq.Page, searchReq.Limit, totalData)
}

// GetPermissionByID implements [PermissionControllerInterface].
func (p *PermissionController) GetPermissionByID(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	permissionID := c.Params("id")
	if permissionID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(permissionID)
	permissionModel, err := p.permissionUsecase.GetPermissionByID(ctx, id)
	if err != nil {
		logger.FailIfError(2, err)
		return httpresponse.Error(c, err)
	}

	permissionResp := permission.NewPermissionResponse(permissionModel)
	return httpresponse.Success(c, httpresponse.FetchMessage(p.entity, true), permissionResp)
}

// CreatePermission implements [PermissionControllerInterface].
func (p *PermissionController) CreatePermission(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	req := permission.PermissionRequest{}

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

	reqModel := models.Permission{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := p.permissionUsecase.CreatePermission(ctx, reqModel); err != nil {
		logger.FailIfError(3, err)
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, httpresponse.CreateMessage(p.entity, true), nil)
}

// UpdatePermission implements [PermissionControllerInterface].
func (p *PermissionController) UpdatePermission(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	permissionID := c.Params("id")
	if permissionID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(permissionID)
	req := permission.PermissionRequest{}

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

	reqModel := models.Permission{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := p.permissionUsecase.UpdatePermission(ctx, reqModel); err != nil {
		logger.FailIfError(3, err)
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, httpresponse.UpdateMessage(p.entity, true), nil)
}

// DeletePermission implements [PermissionControllerInterface].
func (p *PermissionController) DeletePermission(c fiber.Ctx) error {
	ctx := c.RequestCtx()

	permissionID := c.Params("id")
	if permissionID == "" {
		logger.FailIfError(1, errors.New("id parameter is required"))
		details := httpresponse.ErrorDetail(errors.New("id parameter is required"))
		return httpresponse.BadRequest(c, details)
	}

	id := conv.StringToUint(permissionID)
	if err := p.permissionUsecase.DeletePermission(ctx, id); err != nil {
		logger.FailIfError(2, err)
		return httpresponse.Error(c, err)
	}

	return httpresponse.Success(c, httpresponse.DeleteMessage(p.entity, true), nil)
}
