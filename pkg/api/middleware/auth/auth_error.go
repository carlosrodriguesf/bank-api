package auth

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
)

var errorMap = map[error]*apierror.ApiError{
	pkgerror.ErrSessionNotFound: apierror.ErrUnauthorized,
	pkgerror.ErrCantGetSession:  apierror.ErrInternal,
}
