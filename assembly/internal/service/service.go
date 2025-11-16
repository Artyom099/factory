package service

import (
	"context"

	"github.com/Artyom099/factory/assembly/internal/model"
)

type IConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type IProducerService interface {
	ProduceOrderAssembled(ctx context.Context, event model.ShipAssembledOutEvent) error
}
