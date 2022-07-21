package account

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"net/http"
)

var (
	errorMap = map[error]*apierror.ApiError{
		pkgerror.ErrDocumentAlreadyExists: apierror.NewApiError(http.StatusBadRequest, pkgerror.ErrDocumentAlreadyExists.Error(), nil),
		pkgerror.ErrCantCreateAccount:     apierror.NewApiError(http.StatusInternalServerError, pkgerror.ErrCantCreateAccount.Error(), nil),
		pkgerror.ErrCantListAccounts:      apierror.NewApiError(http.StatusInternalServerError, pkgerror.ErrCantListAccounts.Error(), nil),
	}
)
