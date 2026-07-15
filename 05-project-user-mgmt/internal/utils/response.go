package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	ErrCodeInternal   ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeNotFound   ErrorCode = "NOT_FOUND"
	ErrCodeConflict   ErrorCode = "CONFLICT"
	ErrCodeBadRequest ErrorCode = "BAD_REQUEST"
)

type ErrorCode string

type AppError struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (ae *AppError) Error() string {
	return ""
}

func WrapError(err error, msg string, code ErrorCode) error {
	return &AppError{
		Message: msg,
		Code:    code,
		Err:     err,
	}
}

func NewError(message string, code ErrorCode) error {
	return &AppError{
		Message: message,
		Code:    code,
	}
}

func ResponseValidator(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusBadRequest, data)
}

func ResponseError(ctx *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		status := httpStatusFromCode(appErr.Code)
		response := gin.H{
			"error": appErr.Message,
			"code":  appErr.Code,
		}

		if appErr.Err != nil {
			response["detail"] = appErr.Err.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrCodeInternal,
	})
}

func ResponseSuccess(ctx *gin.Context, status int, data any) {
	ctx.JSON(status, gin.H{
		"status": "success",
		"data":   data,
	})
}

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func ResponseStatusCode(ctx *gin.Context, status int) {
	ctx.Status(status)
}
