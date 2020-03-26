package errors

import "errors"

var (
	ErrNotRegistered     = errors.New("username not registered")
	ErrIncorrectPassword = errors.New("password incorrect")
	ErrNotVerified       = errors.New("user not verified")
	ErrAlreadyVerified   = errors.New("user already verified")
)
