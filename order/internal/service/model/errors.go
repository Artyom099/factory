package model

import "errors"

var ErrOrderNotFound = errors.New("order not found")

var ErrConflict = errors.New("conflict")

var ErrNotAllPartsExist = errors.New("not all parts exist")

var ErrListPartsError = errors.New("list parts error")

var ErrOrderAlreadyPaid = errors.New("order already paid")

var ErrOrderCancelled = errors.New("order cancelled")

var ErrInPaymeentService = errors.New("order cancelled")

var ErrUpdateOrder = errors.New("update order error")
