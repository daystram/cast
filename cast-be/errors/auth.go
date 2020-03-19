package errors

import "errors"

var (
	ErrNotRegistered     = errors.New("username not registered")
	ErrIncorrectPassword = errors.New("password incorrect")
)
