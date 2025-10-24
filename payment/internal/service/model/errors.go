package model

import "errors"

var (
	ErrInvalidPaymentMethod = errors.New("invalid payment method")

	ErrInternalError = errors.New("internal server error")
)
