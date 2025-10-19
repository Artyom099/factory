package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Cancel(ctx context.Context, dto model.OrderCancelServiceRequestDto) (model.OrderCancelServiceResponseDto, error) {
	order, err := s.orderRepository.Get(ctx, dto.OrderUUID)
	if err != nil {
		return model.OrderCancelServiceResponseDto{}, model.ErrOrderNotFound
	}

	if model.OrderStatus(order.Status) == model.OrderStatusPAID {
		return model.OrderCancelServiceResponseDto{}, model.ErrConflict
	}

	if model.OrderStatus(order.Status) == model.OrderStatusPENDINGPAYMENT {
		_, err := s.orderRepository.Cancel(ctx, dto.OrderUUID)
		if err != nil {
			return model.OrderCancelServiceResponseDto{}, err
		}
	}

	return model.OrderCancelServiceResponseDto{}, nil
}
