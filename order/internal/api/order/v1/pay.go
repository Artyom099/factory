package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.OrderPayRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	orderUUID, err := uuid.Parse(params.OrderUUID.String())
	if err != nil {
		logger.Error(ctx, "invalid order_uuid", zap.String("order_uuid", orderUUID.String()))
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("invalid order_uuid: %v", orderUUID),
		}, nil
	}

	if err := req.Validate(); err != nil {
		logger.Error(ctx, "order not found", zap.String("order_uuid", orderUUID.String()))
		return &orderV1.BadRequestError{
			Code:    400,
			Message: fmt.Sprintf("validation error: %v", err),
		}, nil
	}

	transactionUUID, err := a.orderService.Pay(ctx, orderUUID.String(), model.OrderPaymentMethod(req.PaymentMethod))
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			logger.Error(ctx, "order not found", zap.String("order_uuid", orderUUID.String()))
			return &orderV1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order %s not found", orderUUID.String()),
			}, nil
		}

		if errors.Is(err, model.ErrOrderAlreadyPaid) {
			logger.Error(ctx, "order already paid", zap.String("order_uuid", orderUUID.String()))
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s already paid", orderUUID.String()),
			}, nil
		}

		if errors.Is(err, model.ErrOrderCancelled) {
			logger.Error(ctx, "order cancelled", zap.String("order_uuid", orderUUID.String()))
			return &orderV1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s cancelled", orderUUID.String()),
			}, nil
		}

		logger.Error(ctx, "internal server error", zap.String("order_uuid", orderUUID.String()), zap.Error(err))
		return &orderV1.InternalServerError{
			Code:    500,
			Message: "Internal Server Error",
		}, err
	}

	return &orderV1.OrderPayResponse{TransactionUUID: transactionUUID}, nil
}
