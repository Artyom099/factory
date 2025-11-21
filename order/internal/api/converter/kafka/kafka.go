package kafka

import "github.com/Artyom099/factory/order/internal/service/model"

type IOrderAssembledDecoder interface {
	Decode(data []byte) (model.OrderAssembledInEvent, error)
}
