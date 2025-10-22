package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

func (r *repository) Cancel(ctx context.Context, orderUuid string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	order.Status = model.OrderStatusCANCELLED

	return nil
}
