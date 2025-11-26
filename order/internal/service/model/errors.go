package model

import "errors"

var (
	ErrOrderNotFound = errors.New("order not found")

	ErrConflict = errors.New("conflict")

	ErrNotAllPartsExist = errors.New("not all parts exist")

	ErrListPartsError = errors.New("list parts error")

	ErrOrderAlreadyPaid = errors.New("order already paid")

	ErrOrderCancelled = errors.New("order cancelled")

	ErrOrderAssembled = errors.New("order assembled")

	ErrOrderPendingPayment = errors.New("order pending payment")

	ErrInPaymeentService = errors.New("payment service error")

	ErrSendOderPaidMessageToKafka = errors.New("send order paid message to kafka error")
)
