package auth

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"net/http"
)

var errorMap = map[error]*apierror.ApiError{
	pkgerror.ErrCantAuth:           apierror.NewApiError(http.StatusInternalServerError, pkgerror.ErrCantAuth.Error(), nil),
	pkgerror.ErrInvalidCredentials: apierror.NewApiError(http.StatusUnauthorized, pkgerror.ErrInvalidCredentials.Error(), nil),
}
