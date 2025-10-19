package service

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
)

type IOrderService interface {
	Create(ctx context.Context, dto model.OrderCreateServiceRequestDto) (model.OrderCreateServiceResponseDto, error)
	Get(ctx context.Context, dto model.OrderGetServiceRequestDto) (model.OrderGetServiceResponseDto, error)
	Cancel(ctx context.Context, dto model.OrderCancelServiceRequestDto) (model.OrderCancelServiceResponseDto, error)
	Pay(ctx context.Context, dto model.OrderPayServiceRequestDto) (model.OrderPayServiceResponseDto, error)
}
