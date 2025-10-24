package repository

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
)

type IOrderRepository interface {
	Create(ctx context.Context, dto model.Order) (string, error)
	Get(ctx context.Context, orderUuid string) (model.Order, error)
	Cancel(ctx context.Context, orderUuid string) error
	Update(ctx context.Context, dto model.Order) error
}
