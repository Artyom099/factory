package order_assembled_consumer

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
	orderAssembledConsumer kafka.IConsumer
	orderAssembledDecoder  kafkaConverter.IOrderAssembledDecoder
	telegramService        def.INotificationTelegramService
}

func NewService(
	orderAssembledConsumer kafka.IConsumer,
	orderAssembledDecoder kafkaConverter.IOrderAssembledDecoder,
	telegramService def.INotificationTelegramService,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramService:        telegramService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderAssembledConsumer.Consume(ctx, s.OrderAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
