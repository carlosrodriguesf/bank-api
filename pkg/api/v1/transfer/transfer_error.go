package transfer

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"net/http"
)

var errorMap = map[error]*apierror.ApiError{
	pkgerror.ErrCantCreateTransfer:            apierror.NewApiError(http.StatusInternalServerError, pkgerror.ErrCantCreateTransfer.Error(), nil),
	pkgerror.ErrOriginAccountTransferNotFound: apierror.NewApiError(http.StatusBadRequest, pkgerror.ErrOriginAccountTransferNotFound.Error(), nil),
	pkgerror.ErrTargetAccountTransferNotFound: apierror.NewApiError(http.StatusBadRequest, pkgerror.ErrTargetAccountTransferNotFound.Error(), nil),
	pkgerror.ErrInsufficientFunds:             apierror.NewApiError(http.StatusBadRequest, pkgerror.ErrInsufficientFunds.Error(), nil),
}
