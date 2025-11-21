package kafka

import "github.com/Artyom099/factory/notification/internal/model"

type IOrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidInEvent, error)
}

type IOrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledInEvent, error)
}
