package error

import (
	"fmt"
	"net/http"
)

var (
	ErrUnauthorized   = NewApiError(http.StatusUnauthorized, "api.access-denied", nil)
	ErrAccessDenied   = NewApiError(http.StatusForbidden, "api.access-denied", nil)
	ErrInternal       = NewApiError(http.StatusInternalServerError, "api.unknown", nil)
	ErrInvalidPayload = NewApiError(http.StatusBadRequest, "api.invalid_payload", nil)
)

type ApiError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Detail  interface{} `json:"detail,omitempty"`
}

func NewApiError(httpCode int, message string, detail interface{}) *ApiError {
	return &ApiError{
		Code:    httpCode,
		Message: message,
		Detail:  detail,
	}
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("%d %s", e.Code, e.Message)
}
