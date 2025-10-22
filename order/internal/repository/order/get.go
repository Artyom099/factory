package order

import (
	"context"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (model.RepoOrder, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[orderUuid]
	if !ok {
		return model.RepoOrder{}, model.ErrOrderNotFound
	}

	return *order, nil
}
