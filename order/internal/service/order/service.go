package order

import (
	"github.com/Artyom099/factory/order/internal/repository"
	def "github.com/Artyom099/factory/order/internal/service"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

var _ def.IOrderService = (*service)(nil)

type service struct {
	orderRepository repository.IOrderRepository
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewService(
	orderRepository repository.IOrderRepository,
	inventoryClient inventoryV1.InventoryServiceClient,
	paymentClient paymentV1.PaymentServiceClient,
) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
