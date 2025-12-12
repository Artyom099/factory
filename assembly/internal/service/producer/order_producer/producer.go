package order_producer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Artyom099/factory/assembly/internal/model"
	"github.com/Artyom099/factory/platform/pkg/kafka"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
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

func (p *service) ProduceOrderAssembled(ctx context.Context, event model.OrderAssembledOutEvent) error {
	ctx, span := tracing.StartSpan(ctx, "assembly.order_assembled_producer",
		trace.WithAttributes(
			attribute.String("event.uuid", event.EventUUID),
			attribute.String("order.uuid", event.OrderUUID),
			attribute.String("user.uuid", event.UserUUID),
		),
	)
	defer span.End()

	msg := &eventsV1.OrderAssembledEvent{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	err = p.orderAssembledProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "failed to publish OrderAssembled", zap.Error(err))
		return err
	}

	return nil
}
