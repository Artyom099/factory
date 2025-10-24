package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/order/internal/service/model"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("Invalid order_uuid: %v", err),
		}, nil
	}

	err = a.orderService.Cancel(ctx, orderUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", orderUUID.String()),
			}, nil
		}
		if errors.Is(err, model.ErrConflict) {
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s paid or already cancelled", orderUUID.String()),
			}, nil
		}

		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, err
	}

	return nil, nil
}
