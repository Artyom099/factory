package service

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/service/model"
)

type IPaymentService interface {
	PayOrder(ctx context.Context, dto model.Payment) (string, error)
}
