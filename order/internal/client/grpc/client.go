package grpc

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
)

type IInventoryClient interface {
	ListParts(ctx context.Context, filter model.ListPartsFilter) (model.ListPartsResponseDto, error)
}

type IPaymentClient interface {
	PayOrder(ctx context.Context, paymentMethod model.OrderPaymentMethod, orderUUID, userUUID string) (string, error)
}
