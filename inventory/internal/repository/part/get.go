package part

import (
	"context"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, dto repoModel.PartGetRepoRequest) (repoModel.PartGetRepoResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[dto.Uuid]
	if !ok {
		return repoModel.PartGetRepoResponse{}, repoModel.ErrPartNotFound
	}

	return repoModel.PartGetRepoResponse{Part: part}, nil
}
