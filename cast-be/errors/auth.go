package errors

import "errors"

// Error returns
var (
	ErrNotRegistered     = errors.New("username not registered")
	ErrIncorrectPassword = errors.New("password incorrect")
	ErrNotVerified       = errors.New("user not verified")
	ErrAlreadyVerified   = errors.New("user already verified")
)
