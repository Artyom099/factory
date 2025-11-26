package order

import (
	"context"
	"errors"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Assemble(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	if order.Status == model.OrderStatusPENDINGPAYMENT {
		return model.ErrOrderPendingPayment
	}
	if order.Status == model.OrderStatusASSEMBLED {
		return model.ErrOrderAssembled
	}
	if order.Status == model.OrderStatusCANCELLED {
		return model.ErrOrderCancelled
	}

	order.Status = model.OrderStatusASSEMBLED
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	return nil
}
