package repository

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/service/model"
)

type IPaymentRepository interface {
	PayOrder(ctx context.Context, dto model.Payment) (string, error)
}
