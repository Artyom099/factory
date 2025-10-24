package order

import (
	"context"
	"errors"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.OrderPaymentMethod) (string, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return "", model.ErrOrderNotFound
		}
		return "", model.ErrInternalError
	}

	if order.Status == model.OrderStatusPAID {
		return "", model.ErrOrderAlreadyPaid
	}

	if order.Status == model.OrderStatusCANCELLED {
		return "", model.ErrOrderCancelled
	}

	transactionUUID, err := s.paymentClient.PayOrder(ctx, paymentMethod, orderUUID, order.UserUUID)
	if err != nil {
		return "", model.ErrInPaymeentService
	}

	order.Status = model.OrderStatusPAID
	order.TransactionUUID = transactionUUID
	order.PaymentMethod = paymentMethod

	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return "", model.ErrUpdateOrder
	}

	return transactionUUID, nil
}
