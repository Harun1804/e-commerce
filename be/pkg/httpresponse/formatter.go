package httpresponse

import "github.com/gofiber/fiber/v3"

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
	Data       interface{} `json:"data,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// JSON sends a formatted response. If `pagination` is nil it will be omitted.
func JSON(c fiber.Ctx, code int, status bool, message string, data interface{}, pagination *Pagination) error {
	if code < 100 || code > 599 {
		code = fiber.StatusOK
	}

	resp := Response{
		Code:       code,
		Status:     status,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	}

	return c.Status(code).JSON(resp)
}

// Success sends a 200 OK response with data.
func Success(c fiber.Ctx, code int, message string, data interface{}) error {
	return JSON(c, code, true, message, data, nil)
}

// SuccessWithPagination sends a 200 OK response including pagination info.
func SuccessWithPagination(c fiber.Ctx, code int, message string, data interface{}, page, limit int, totalData int64) error {
	p := generatePagination(page, limit, totalData)
	return JSON(c, code, true, message, data, p)
}

func generatePagination(page, limit int, totalData int64) *Pagination {
	return &Pagination{
		Page:      page,
		Size:      limit,
		TotalData: totalData,
		TotalPage: (totalData + int64(limit) - 1) / int64(limit),
	}
}

// Error sends an error response using the provided HTTP status code.
func Error(c fiber.Ctx, code int, message string, data *interface{}) error {
	if code < 100 || code > 599 {
		code = fiber.StatusBadRequest
	}

	if data != nil {
		return JSON(c, code, false, message, data, nil)
	}

	return JSON(c, code, false, message, nil, nil)
}

func BadRequest(c fiber.Ctx, message string, data *interface{}) error {
	return Error(c, fiber.StatusBadRequest, message, data)
}

func NotFound(c fiber.Ctx, message string) error {
	return Error(c, fiber.StatusNotFound, message, nil)
}

func ValidationError(c fiber.Ctx, data *interface{}) error {
	return Error(c, fiber.StatusUnprocessableEntity, "Validation failed", data)
}

func InternalServerError(c fiber.Ctx, message string, data *interface{}) error {
	return Error(c, fiber.StatusInternalServerError, message, data)
}
