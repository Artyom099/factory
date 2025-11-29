package service

import (
	"context"

	"github.com/Artyom099/factory/assembly/internal/model"
)

type IAssemblyService interface {
	Assembly(ctx context.Context, dto model.OrderPaidInEvent) error
}

type IAssemblyConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type IAssemblyProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledOutEvent) error
}
