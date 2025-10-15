package service

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/model"
)

type IPaymentService interface {
	PayOrder(ctx context.Context, dto model.PayOrderServiceRequestDto) (string, error)
}
