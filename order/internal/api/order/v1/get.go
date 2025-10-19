package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/order/internal/api/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Invalid order_uuid: %v", err),
		}, nil
	}

	res, err := a.orderService.Get(ctx, converter.OrderGetApiRequestDtoToOrderGetServiceRequestDto(params))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
			}, nil
		}
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, err
	}

	return converter.OrderGetServiceResponseDtoToOrderGetApiResponseDto(res), nil
}
