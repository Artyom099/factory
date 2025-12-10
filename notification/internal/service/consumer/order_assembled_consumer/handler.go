package order_assembled_consumer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) OrderAssembledHandler(ctx context.Context, msg consumer.Message) error {
	ctx, orderAssembledSpan := tracing.StartSpan(ctx, "notification.order_assembled_consumer",
		trace.WithAttributes(
			attribute.String("topic", msg.Topic),
		),
	)

	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		orderAssembledSpan.RecordError(err)
		orderAssembledSpan.End()
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err), zap.Any("message: ", msg))
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

	ctx, notificationSpan := tracing.StartSpan(ctx, "notification.send_order_assembled_notification",
		trace.WithAttributes(
			attribute.String("order.uuid", event.OrderUUID),
		),
	)
	defer notificationSpan.End()

	if err = s.telegramService.SendOrderAssembledNotification(ctx, event); err != nil {
		notificationSpan.RecordError(err)
		return err
	}

	return nil
}
