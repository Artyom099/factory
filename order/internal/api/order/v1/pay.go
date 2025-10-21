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

func (a *api) PayOrder(ctx context.Context, req *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	if _, err := uuid.Parse(params.OrderUUID.String()); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("invalid order_uuid: %v", err),
		}, nil
	}

	if err := req.Validate(); err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, converter.OrderPayApiRequestDtoToOrderPayServiceRequestDto(params, req))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", params.OrderUUID.String()),
			}, nil
		}
		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s already paid", params.OrderUUID.String()),
			}, nil
		}
		if errors.Is(err, model.ErrOrderCancelled) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s cancelled", params.OrderUUID.String()),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, err
	}

	return &orderV1.OrderPayResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
