package repository

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

type IOrderRepository interface {
	Create(ctx context.Context, dto model.RepoOrder) (string, error)
	Get(ctx context.Context, orderUuid string) (model.RepoOrder, error)
	Cancel(ctx context.Context, orderUuid string) error
	Update(ctx context.Context, dto model.RepoOrder) error
}
