package model

import "errors"

var (
	ErrSessionNotFound = errors.New("session not found")

	ErrInvalidPassword = errors.New("invalid password")
)
