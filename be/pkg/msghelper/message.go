package msghelper

import "fmt"

func AlreadyExistsMessage(entity string, identifier any) string {
	return entity + " with identifier '" + fmt.Sprint(identifier) + "' already exists"
}

func NotFoundMessage(entity string, id uint) string {
	return entity + " with ID " + fmt.Sprint(id) + " not found"
}

func NotFoundMessageByName(entity string, name string) string {
	return entity + " with name '" + name + "' not found"
}

func FetchMessage(entity string, isSuccess bool) string {
	if isSuccess {
		return entity + " retrieved successfully"
	}
	return entity + " retrieval failed"
}

func CreateMessage(entity string, isSuccess bool) string {
	if isSuccess {
		return entity + " created successfully"
	}
	return entity + " creation failed"
}

func UpdateMessage(entity string, isSuccess bool) string {
	if isSuccess {
		return entity + " updated successfully"
	}
	return entity + " update failed"
}

func DeleteMessage(entity string, isSuccess bool) string {
	if isSuccess {
		return entity + " deleted successfully"
	}
	return entity + " deletion failed"
}