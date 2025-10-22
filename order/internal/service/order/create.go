package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *service) Create(ctx context.Context, dto model.Order) (model.Order, error) {
	parts, err := s.inventoryClient.ListParts(ctx, model.ListPartsFilter{
		Uuids: dto.PartUuids,
	})
	if err != nil {
		return model.Order{}, model.ErrListPartsError
	}

	if len(parts) != len(dto.PartUuids) {
		return model.Order{}, model.ErrNotAllPartsExist
	}

	var totalPrice float32
	for _, part := range parts {
		totalPrice += float32(part.Price)
	}

	orderUuid, err := s.orderRepository.Create(ctx, converter.ModelToRepoOrder(dto, totalPrice))
	if err != nil {
		return model.Order{}, err
	}

	return model.Order{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
