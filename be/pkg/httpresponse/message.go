package httpresponse

import "fmt"

func AlreadyExistsMessage(entity string, identifier any) string {
	return entity + " with identifier '" + fmt.Sprint(identifier) + "' already exists"
}

func NotFoundMessage(entity, field string, identifier any) string {
	return identifierNotFoundMessage(entity, field, identifier)
}

func InternalServerErrorMessage(message string) string {
	if message == "" {
		return "Internal Server Error"
	}

	return "Internal Server Error: " + message
}

func FetchMessage(entity string, isSuccess bool) string {
	return statusMessage(entity, "retrieved successfully", "retrieval failed", isSuccess)
}

func CreateMessage(entity string, isSuccess bool) string {
	return statusMessage(entity, "created successfully", "creation failed", isSuccess)
}

func UpdateMessage(entity string, isSuccess bool) string {
	return statusMessage(entity, "updated successfully", "update failed", isSuccess)
}

func DeleteMessage(entity string, isSuccess bool) string {
	return statusMessage(entity, "deleted successfully", "deletion failed", isSuccess)
}

func identifierNotFoundMessage(entity, field string, value any) string {
	return entity + " with " + field + " '" + fmt.Sprint(value) + "' not found"
}

func statusMessage(entity, successMessage, failedMessage string, isSuccess bool) string {
	if isSuccess {
		return entity + " " + successMessage
	}

	return entity + " " + failedMessage
}