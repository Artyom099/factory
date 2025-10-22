package part

import (
	"context"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (repoModel.RepoPart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[uuid]
	if !ok {
		return repoModel.RepoPart{}, repoModel.ErrPartNotFound
	}

	return repoModel.RepoPart(part), nil
}
