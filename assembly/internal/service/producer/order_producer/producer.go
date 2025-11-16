package order_producer

import (
	"context"

	"github.com/gogo/protobuf/proto"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/assembly/internal/model"
	"github.com/Artyom099/factory/platform/pkg/kafka"
	"github.com/Artyom099/factory/platform/pkg/logger"
	eventsV1 "github.com/Artyom099/factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderAssembledProducer kafka.IProducer
}

func NewService(orderAssembledProducer kafka.IProducer) *service {
	return &service{
		orderAssembledProducer: orderAssembledProducer,
	}
}

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.ShipAssembledOutEvent) error {
	msg := &eventsV1.ShipAssembledEvent{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderAssembled", zap.Error(err))
		return err
	}

	return nil
}
