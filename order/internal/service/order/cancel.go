package order

import (
	"context"
	"errors"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	if order.Status == model.OrderStatusPAID || order.Status == model.OrderStatusCANCELLED {
		return model.ErrConflict
	}

	if order.Status == model.OrderStatusPENDINGPAYMENT {
		err := s.orderRepository.Cancel(ctx, orderUUID)
		if err != nil {
			if errors.Is(err, repoModel.ErrOrderNotFound) {
				return model.ErrOrderNotFound
			}
			return err
		}
	}

	return nil
}
