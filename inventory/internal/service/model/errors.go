package model

import "errors"

var (
	ErrPartNotFound = errors.New("part not found")

	ErrInternalError = errors.New("internal error")
)
