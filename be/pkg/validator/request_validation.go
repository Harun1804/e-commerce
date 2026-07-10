package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(data interface{}) error {
	var errorMessages []string
	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errorMessages = append(errorMessages, err.Field()+" is required")
			case "email":
				errorMessages = append(errorMessages, err.Field()+" must be a valid email")
			case "min":
				errorMessages = append(errorMessages, err.Field()+" must be at least "+err.Param()+" characters long")
			case "max":
				errorMessages = append(errorMessages, err.Field()+" must be at most "+err.Param()+" characters long")
			}
		}
		return errors.New("Validasi gagal : " + joinMessages(errorMessages))
	}
	return nil
}

func joinMessages(messages []string) string {
	result := ""
	for i, msg := range messages {
		if i > 0 {
			result += ", "
		}
		result += msg
	}
	return result
}
