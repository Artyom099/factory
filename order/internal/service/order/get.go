package order

import (
	"context"
	"errors"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Get(ctx context.Context, dto servModel.OrderGetServiceRequestDto) (servModel.OrderGetServiceResponseDto, error) {
	res, err := s.orderRepository.Get(ctx, dto.OrderUUID)
	if err != nil {
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return servModel.OrderGetServiceResponseDto{}, servModel.ErrOrderNotFound
		}
		return servModel.OrderGetServiceResponseDto{}, err
	}

	return converter.OrderGetRepoResponseDtoToOrderGetServiceResponseDto(res), nil
}
