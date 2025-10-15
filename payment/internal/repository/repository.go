package repository

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/repository/model"
)

type IPaymentRepository interface {
	PayOrder(ctx context.Context, dto model.PayOrderRepoRequestDto) (string, error)
}
