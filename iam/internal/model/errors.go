package model

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")

	ErrUserAlreadyExists = errors.New("user already exists")

	ErrUserNotFound = errors.New("user not found")

	ErrInvalidPassword = errors.New("invalid password")
)
