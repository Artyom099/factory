package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (s *service) Create(ctx context.Context, dto model.OrderCreateServiceRequestDto) (model.OrderCreateServiceResponseDto, error) {
	listPartsReq := inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: dto.PartUuids,
		},
	}
	parts, err := s.inventoryClient.ListParts(ctx, &listPartsReq)
	if err != nil {
		return model.OrderCreateServiceResponseDto{}, model.ErrListPartsError
	}

	if len(parts.GetParts()) != len(dto.PartUuids) {
		return model.OrderCreateServiceResponseDto{}, model.ErrNotAllPartsExist
	}

	var totalPrice float32
	for _, part := range parts.GetParts() {
		totalPrice += float32(part.GetPrice())
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
