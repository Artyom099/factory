package order

import (
	"context"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Update(ctx context.Context, dto model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[dto.OrderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	order.TransactionUUID = dto.TransactionUUID
	order.Status = repoModel.OrderStatus(dto.Status)
	order.PaymentMethod = repoModel.OrderPaymentMethod(dto.PaymentMethod)

	return nil
}
