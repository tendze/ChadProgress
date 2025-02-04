package service

import (
	"errors"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrFieldIsTooLong     = errors.New("field is too long")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateKey       = errors.New("duplicate key value violates unique constraint")
	ErrInvalidRoleRequest = errors.New("invalid request due to role")
	ErrUserNotFound       = errors.New("user not found")
	ErrClientNotFound     = errors.New("clients profile not found")
	ErrTrainerNotFound    = errors.New("trainer profile not found")
	ErrNotActiveTrainer   = errors.New("not active trainer")
)
