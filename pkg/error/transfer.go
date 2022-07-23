package errors

import "errors"

var (
	ErrCantCreateTransfer            = errors.New("transfer.cant-create-transfer")
	ErrOriginAccountTransferNotFound = errors.New("transfer.origin-not-found")
	ErrTargetAccountTransferNotFound = errors.New("transfer.target-not-found")
)
