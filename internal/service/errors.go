package service

import "errors"

var (
	ErrUserExists         = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenGeneration    = errors.New("failed to generate token")
	ErrTaskNotFound       = errors.New("task not found")
	ErrUnauthorized       = errors.New("unauthorized access")
)
