package service

import (
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrFieldIsTooLong     = errors.New("field is too long")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateKey       = errors.New("duplicate key value violates unique constraint")
)
