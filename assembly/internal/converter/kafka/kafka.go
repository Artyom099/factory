package kafka

import "github.com/Artyom099/factory/assembly/internal/model"

type IOrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidInEvent, error)
}
