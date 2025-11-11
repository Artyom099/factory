package v1

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/Artyom099/factory/order/internal/api/converter"
	"github.com/Artyom099/factory/platform/pkg/logger"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.OrderCreateRequest) (orderV1.CreateOrderRes, error) {
	if err := req.Validate(); err != nil {
		logger.Error(ctx, "validation error", zap.Error(err))
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	res, err := a.orderService.Create(ctx, converter.ToModelOrder(req))
	if err != nil {
		logger.Error(ctx, "internal server error", zap.Error(err))
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, nil
	}

	return &orderV1.OrderCreateResponse{
		OrderUUID:  res.OrderUUID,
		TotalPrice: res.TotalPrice,
	}, nil
}
