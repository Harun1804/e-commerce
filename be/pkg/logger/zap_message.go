package logger

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func FailIfError(order int, err error, fields ...zap.Field) error {
	if err == nil {
		return nil
	}

	zap.L().Error(wrapCallerMessage(order), append(fields, zap.Error(err))...)
	return err
}

func FailIfRowsAffectedZero(order int, rowsAffected int64, fields ...zap.Field) error {
	if rowsAffected != 0 {
		return nil
	}

	zap.L().Error(wrapCallerMessage(order), fields...)
	return gorm.ErrRecordNotFound
}

func wrapCallerMessage(order int) string {
	fileName, methodName := callerContext()
	return fmt.Sprintf("[%s] %s - %d", fileName, methodName, order)
}

func callerContext() (string, string) {
	callerPC, callerFile, _, ok := runtime.Caller(3)
	if !ok {
		return "unknown.go", "unknown"
	}

	methodName := "unknown"
	if callerFunc := runtime.FuncForPC(callerPC); callerFunc != nil {
		methodName = shortFunctionName(callerFunc.Name())
	}

	return path.Base(callerFile), methodName
}

func shortFunctionName(functionName string) string {
	functionName = strings.TrimSuffix(functionName, "-fm")
	if lastSlash := strings.LastIndex(functionName, "/"); lastSlash >= 0 {
		functionName = functionName[lastSlash+1:]
	}

	if lastDot := strings.LastIndex(functionName, "."); lastDot >= 0 {
		functionName = functionName[lastDot+1:]
	}

	return functionName
}
