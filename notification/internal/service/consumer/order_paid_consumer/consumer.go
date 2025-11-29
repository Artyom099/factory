package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Artyom099/factory/notification/internal/converter/kafka"
	def "github.com/Artyom099/factory/notification/internal/service"
	"github.com/Artyom099/factory/platform/pkg/kafka"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

var _ def.INotificationConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.IConsumer
	orderPaidDecoder  kafkaConverter.IOrderPaidDecoder
	telegramService   def.INotificationTelegramService
}

func NewService(
	orderPaidConsumer kafka.IConsumer,
	orderPaidDecoder kafkaConverter.IOrderPaidDecoder,
	telegramService def.INotificationTelegramService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramService:   telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderPaidHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
