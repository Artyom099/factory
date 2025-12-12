package order_producer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/kafka"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
	eventsV1 "github.com/Artyom099/factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderPaidProducer kafka.IProducer
}

func NewService(orderPaidProducer kafka.IProducer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidOutEvent) error {
	ctx, span := tracing.StartSpan(ctx, "order.order_paid_producer",
		trace.WithAttributes(
			attribute.String("event.uuid", event.EventUUID),
			attribute.String("order.uuid", event.OrderUUID),
			attribute.String("user.uuid", event.UserUUID),
		),
	)
	defer span.End()

	msg := &eventsV1.OrderPaidEvent{
		EventUuid: event.EventUUID,
		OrderUuid: event.OrderUUID,
		UserUuid:  event.UserUUID,
		PaymentMethod: eventsV1.PaymentMethod(
			eventsV1.PaymentMethod_value[event.PaymentMethod],
		),
		TransactionUuid: event.TransactionUUID,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "failed to marshal OrderPaid event", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		span.RecordError(err)
		logger.Error(ctx, "failed to publish OrderPaid event", zap.Error(err))
		return err
	}

	return nil
}
