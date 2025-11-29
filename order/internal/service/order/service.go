package order

import (
	"github.com/Artyom099/factory/order/internal/client/grpc"
	"github.com/Artyom099/factory/order/internal/repository"
	def "github.com/Artyom099/factory/order/internal/service"
)

var _ def.IOrderService = (*service)(nil)

type service struct {
	orderRepository repository.IOrderRepository

	inventoryClient grpc.IInventoryClient
	paymentClient   grpc.IPaymentClient

	orderProducerService def.IOrderProducerService
	orderConsumerService def.IOrderConsumerService
}

func NewService(
	orderRepository repository.IOrderRepository,
	inventoryClient grpc.IInventoryClient,
	paymentClient grpc.IPaymentClient,
	orderProducerService def.IOrderProducerService,
	orderConsumerService def.IOrderConsumerService,
) *service {
	return &service{
		orderRepository:      orderRepository,
		inventoryClient:      inventoryClient,
		paymentClient:        paymentClient,
		orderProducerService: orderProducerService,
		orderConsumerService: orderConsumerService,
	}
}
