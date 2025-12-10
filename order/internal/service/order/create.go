package order

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Create(ctx context.Context, dto model.Order) (model.Order, error) {
	ctx, listPartsSpan := tracing.StartSpan(ctx, "order.call_inventory_list_parts",
		trace.WithAttributes(
			attribute.String("user.uuid", dto.UserUUID),
			attribute.StringSlice("part.uuids", dto.PartUuids),
		),
	)

	parts, err := s.inventoryClient.ListParts(ctx, model.ListPartsFilter{
		Uuids: dto.PartUuids,
	})
	if err != nil {
		listPartsSpan.RecordError(err)
		listPartsSpan.End()
		return model.Order{}, model.ErrListPartsError
	}

	if len(parts) != len(dto.PartUuids) {
		listPartsSpan.RecordError(err)
		listPartsSpan.End()
		return model.Order{}, model.ErrNotAllPartsExist
	}

	listPartsSpan.End()

	ctx, createOrderSpan := tracing.StartSpan(ctx, "order.create_order",
		trace.WithAttributes(
			attribute.String("user.uuid", dto.UserUUID),
			attribute.StringSlice("part.uuids", dto.PartUuids),
		),
	)
	defer createOrderSpan.End()

	var totalPrice float32
	for _, part := range parts {
		totalPrice += float32(part.Price)
	}

	dto.TotalPrice = totalPrice
	orderUuid, err := s.orderRepository.Create(ctx, dto)
	if err != nil {
		createOrderSpan.RecordError(err)
		return model.Order{}, err
	}

	createOrderSpan.SetAttributes(
		attribute.String("order.uuid", orderUuid),
		attribute.Float64("total.price", float64(totalPrice)),
	)

	return model.Order{
		OrderUUID:  orderUuid,
		TotalPrice: totalPrice,
	}, nil
}
