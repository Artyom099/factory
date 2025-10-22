package repository

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/repository/model"
)

type IPartRepository interface {
	Get(ctx context.Context, uuid string) (model.RepoPart, error)
	List(ctx context.Context, dto model.RepoPartFilter) ([]model.RepoPart, error)
	Create(ctx context.Context, dto model.RepoPart) (string, error)
}
