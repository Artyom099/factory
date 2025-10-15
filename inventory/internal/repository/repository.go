package repository

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/repository/model"
)

type IPartRepository interface {
	Get(ctx context.Context, dto model.PartGetRepoRequest) (model.PartGetRepoResponse, error)
	List(ctx context.Context, dto model.PartListRepoRequest) (model.PartListRepoResponse, error)
	Create(ctx context.Context, dto model.PartCreateRepoRequest) (string, error)
}
