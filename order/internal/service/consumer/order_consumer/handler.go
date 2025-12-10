package order_consumer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) OrderHandler(ctx context.Context, msg consumer.Message) error {
	ctx, orderAssembledSpan := tracing.StartSpan(ctx, "order.order_assembled_consumer",
		trace.WithAttributes(
			attribute.String("topic", msg.Topic),
		),
	)

	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		orderAssembledSpan.RecordError(err)
		orderAssembledSpan.End()
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	orderAssembledSpan.SetAttributes(
		attribute.String("order.uuid", event.OrderUUID),
	)
	orderAssembledSpan.End()

	ctx, assembleSpan := tracing.StartSpan(ctx, "order.assemble_order",
		trace.WithAttributes(
			attribute.String("order.uuid", event.OrderUUID),
		),
	)
	defer assembleSpan.End()

	err = s.orderService.Assemble(ctx, event.OrderUUID)
	if err != nil {
		assembleSpan.RecordError(err)
		logger.Error(ctx, "Failed to update order status to assembled", zap.Error(err))
		return err
	}

	return nil
}
