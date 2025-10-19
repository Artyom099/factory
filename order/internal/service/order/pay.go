package order

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Pay(ctx context.Context, dto model.OrderPayServiceRequestDto) (model.OrderPayServiceResponseDto, error) {
	order, err := s.orderRepository.Get(ctx, dto.OrderUUID)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderNotFound
	}

	if model.OrderStatus(order.Status) == model.OrderStatusPAID {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderAlreadyPaid
	}

	if model.OrderStatus(order.Status) == model.OrderStatusCANCELLED {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderCancelled
	}

	userUuid := uuid.New().String()
	transactionUUID, err := s.paymentClient.PayOrder(ctx, dto.PaymentMethod, dto.OrderUUID, userUuid)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrInPaymeentService
	}

	updateOrderDto := repoModel.OrderUpdateRepoRequestDto{
		OrderUUID:       dto.OrderUUID,
		Status:          repoModel.OrderStatusPAID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   repoModel.OrderPaymentMethod(dto.PaymentMethod),
	}
	_, err = s.orderRepository.Update(ctx, updateOrderDto)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrUpdateOrder
	}

	return model.OrderPayServiceResponseDto{
		TransactionUUID: transactionUUID,
	}, nil
}
