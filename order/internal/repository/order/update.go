package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

func (r *repository) Update(ctx context.Context, dto model.RepoOrder) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[dto.OrderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	order.Status = dto.Status
	order.TransactionUUID = dto.TransactionUUID
	order.PaymentMethod = dto.PaymentMethod

	return nil
}
