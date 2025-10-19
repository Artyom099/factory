package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Create(ctx context.Context, dto model.OrderCreateServiceRequestDto) (model.OrderCreateServiceResponseDto, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.ListPartsFilter{
		Uuids: dto.PartUuids,
	})
	if err != nil {
		return model.OrderCreateServiceResponseDto{}, model.ErrListPartsError
	}

	if len(parts.Parts) != len(dto.PartUuids) {
		return model.OrderCreateServiceResponseDto{}, model.ErrNotAllPartsExist
	}

	var totalPrice float32
	for _, part := range parts.Parts {
		totalPrice += float32(part.Price)
	}

	orderUuid, err := s.orderRepository.Create(ctx, converter.OrderCreateServiceRequestDtoToOrderCreateRepoRequestDto(dto, totalPrice))
	if err != nil {
		return model.OrderCreateServiceResponseDto{}, err
	}

	return model.OrderCreateServiceResponseDto{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
