package repository

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

type IOrderRepository interface {
	Create(ctx context.Context, dto model.OrderCreateRepoRequestDto) (string, error)
	Get(ctx context.Context, orderUuid string) (model.OrderGetRepoResponseDto, error)
	Cancel(ctx context.Context, orderUuid string) (model.OrderCancelRepoResponseDto, error)
	Update(ctx context.Context, dto model.OrderUpdateRepoRequestDto) (model.OrderUpdateRepoResponseDto, error)
}
