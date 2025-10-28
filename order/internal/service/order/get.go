package order

import (
	"context"
	"errors"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (servModel.Order, error) {
	res, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return servModel.Order{}, servModel.ErrOrderNotFound
		}
		return servModel.Order{}, err
	}

	return res, nil
}
