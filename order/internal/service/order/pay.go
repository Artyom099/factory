package order

import (
	"context"
	"errors"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Pay(ctx context.Context, dto model.Order) (string, error) {
	order, err := s.orderRepository.Get(ctx, dto.OrderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", err
	}

	if model.OrderStatus(order.Status) == model.OrderStatusPAID {
		return "", model.ErrOrderAlreadyPaid
	}

	if model.OrderStatus(order.Status) == model.OrderStatusCANCELLED {
		return "", model.ErrOrderCancelled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, dto.PaymentMethod, dto.OrderUUID, order.UserUUID)
	if err != nil {
		return "", model.ErrInPaymeentService
	}

	updateOrderDto := repoModel.RepoOrder{
		OrderUUID:       dto.OrderUUID,
		Status:          repoModel.OrderStatusPAID,
		TransactionUUID: transactionUUID,
		PaymentMethod:   repoModel.OrderPaymentMethod(dto.PaymentMethod),
	}
	err = s.orderRepository.Update(ctx, updateOrderDto)
	if err != nil {
		return "", model.ErrUpdateOrder
	}

	return transactionUUID, nil
}
