package service

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
)

type IOrderService interface {
	Create(ctx context.Context, dto model.Order) (model.Order, error)
	Get(ctx context.Context, orderUUID string) (model.Order, error)
	Cancel(ctx context.Context, orderUuid string) error
	Pay(ctx context.Context, orderUUID string, paymentMethod model.OrderPaymentMethod) (string, error)
}
