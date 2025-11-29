package model

import "errors"

var (
	ErrConsumeOderPaidMessageToKafka = errors.New("consume order paid message to kafka error")

	ErrSendOderAssembledMessageToKafka = errors.New("send order assembled message to kafka error")
)
