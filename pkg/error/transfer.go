package errors

import "errors"

var (
	ErrCantCreateTransfer            = errors.New("transfer.cant-create-transfer")
	ErrCantListTransfers             = errors.New("transfer.cant-list-transfer")
	ErrOriginAccountTransferNotFound = errors.New("transfer.origin-not-found")
	ErrTargetAccountTransferNotFound = errors.New("transfer.target-not-found")
)
