package service

import (
	"context"

	"github.com/Artyom099/factory/assembly/internal/model"
)

type IAssemblyConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type IAssemblyProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.ShipAssembledOutEvent) error
}
