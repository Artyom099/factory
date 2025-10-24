package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return model.Order{}, model.ErrOrderNotFound
	}

	return converter.ToModelOrder(*order), nil
}
