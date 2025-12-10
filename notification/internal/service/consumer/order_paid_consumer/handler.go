package order_paid_consumer

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) OrderPaidHandler(ctx context.Context, msg consumer.Message) error {
	ctx, orderPaidSpan := tracing.StartSpan(ctx, "notification.order_paid_consumer",
		trace.WithAttributes(
			attribute.String("topic", msg.Topic),
		),
	)

	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		orderPaidSpan.RecordError(err)
		orderPaidSpan.End()
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err), zap.Any("message: ", msg))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transation_uuid", event.TransactionUUID),
	)

	orderPaidSpan.SetAttributes(
		attribute.String("order.uuid", event.OrderUUID),
	)
	orderPaidSpan.End()

	ctx, notificationSpan := tracing.StartSpan(ctx, "notification.send_order_paid_notification",
		trace.WithAttributes(
			attribute.String("order.uuid", event.OrderUUID),
		),
	)
	defer notificationSpan.End()

	if err = s.telegramService.SendOrderPaidNotification(ctx, event); err != nil {
		notificationSpan.RecordError(err)
		return err
	}

	return nil
}
