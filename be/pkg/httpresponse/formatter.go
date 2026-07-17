package httpresponse

import (
	"errors"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

type Option func(*Response)

type Pagination struct {
	Page      int   `json:"page"`
	TotalData int64 `json:"totalData"`
	Size      int   `json:"size"`
	TotalPage int64 `json:"totalPage"`
}

type Response struct {
	Code       int         `json:"code"`
	Status     bool        `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination *Pagination `json:"pagination"`
}

func WithPagination(page, size int, totalData int64) Option {
	p := newPagination(page, size, totalData)
	return func(resp *Response) {
		resp.Pagination = p
	}
}

// Success sends a 200 OK response with data.
func Success(c fiber.Ctx, message string, data interface{}) error {
	return send(c, fiber.StatusOK, true, message, data)
}

// SuccessWithPagination sends a 200 OK response including pagination info.
func SuccessWithPagination(c fiber.Ctx, message string, data interface{}, page, limit int, totalData int64) error {
	return send(c, fiber.StatusOK, true, message, data, WithPagination(page, limit, totalData))
}

func BadRequest(c fiber.Ctx, data interface{}) error {
	return sendError(c, fiber.StatusBadRequest, "Bad Request", data)
}

func Unauthorized(c fiber.Ctx) error {
	return sendError(c, fiber.StatusUnauthorized, "Unauthorized", nil)
}

func Forbidden(c fiber.Ctx) error {
	return sendError(c, fiber.StatusForbidden, "Forbidden", nil)
}

func NotFound(c fiber.Ctx, entity, field string, identifier any) error {
	return sendError(c, fiber.StatusNotFound, NotFoundMessage(entity, field, identifier), nil)
}

func Conflict(c fiber.Ctx, entity string, identifier any) error {
	return sendError(c, fiber.StatusConflict, AlreadyExistsMessage(entity, identifier), nil)
}

func ValidationError(c fiber.Ctx, data interface{}) error {
	return sendError(c, fiber.StatusUnprocessableEntity, "Validation failed", data)
}

func InternalServerError(c fiber.Ctx, data interface{}) error {
	return sendError(c, fiber.StatusInternalServerError, InternalServerErrorMessage(""), data)
}

func Error(c fiber.Ctx, err error) error {
	details := ErrorDetail(err)

	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(strings.ToLower(err.Error()), "not found") {
		return sendError(c, fiber.StatusNotFound, err.Error(), details)
	}

	return InternalServerError(c, details)
}

func send(c fiber.Ctx, code int, status bool, message string, data interface{}, opts ...Option) error {
	code = sanitizeStatusCode(code, fiber.StatusOK)

	resp := Response{
		Code:       code,
		Status:     status,
		Message:    message,
		Data:       normalizeData(data),
		Pagination: nil,
	}

	for _, opt := range opts {
		if opt != nil {
			opt(&resp)
		}
	}

	return c.Status(code).JSON(resp)
}

// Error sends an error response using the provided HTTP status code.
func sendError(c fiber.Ctx, code int, message string, data interface{}) error {
	code = sanitizeStatusCode(code, fiber.StatusBadRequest)
	return send(c, code, false, message, data)
}

func sanitizeStatusCode(code int, fallback int) int {
	if code < 100 || code > 599 {
		return fallback
	}

	return code
}

func newPagination(page, size int, totalData int64) *Pagination {
	if totalData <= 0 {
		return nil
	}

	if page < 1 {
		page = 1
	}

	if size < 1 {
		size = 10
	}

	if totalData < 0 {
		totalData = 0
	}

	return &Pagination{
		Page:      page,
		Size:      size,
		TotalData: totalData,
		TotalPage: (totalData + int64(size) - 1) / int64(size),
	}
}

func normalizeData(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array, reflect.String:
		if v.Len() == 0 {
			return nil
		}
	case reflect.Struct:
		if v.NumField() == 0 {
			return nil
		}
	}

	return data
}
