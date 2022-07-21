package errors

import "errors"

var (
	ErrCantCreateAccount     = errors.New("account.cant-create-account")
	ErrDocumentAlreadyExists = errors.New("account.document-already-exists")
	ErrCantListAccounts      = errors.New("account.cant-list-accounts")
)
