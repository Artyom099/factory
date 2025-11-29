package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/Artyom099/factory/assembly/internal/converter/kafka"
	def "github.com/Artyom099/factory/assembly/internal/service"
	"github.com/Artyom099/factory/platform/pkg/kafka"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

var _ def.IAssemblyConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.IConsumer
	orderPaidDecoder  kafkaConverter.IOrderPaidDecoder
	assemblyService   def.IAssemblyService
}

func NewService(
	orderPaidConsumer kafka.IConsumer,
	orderPaidDecoder kafkaConverter.IOrderPaidDecoder,
	assemblyService def.IAssemblyService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		assemblyService:   assemblyService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting orderPaidConsumer service")

	err := s.orderPaidConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
