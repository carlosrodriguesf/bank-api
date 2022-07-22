package errors

import "errors"

var (
	ErrCantAuth           = errors.New("auth.cant-auth")
	ErrInvalidCredentials = errors.New("auth.invalid-credentials")
	ErrCantGetSession     = errors.New("auth.cant-get-session")
	ErrSessionNotFound    = errors.New("auth.session-not-found")
)
