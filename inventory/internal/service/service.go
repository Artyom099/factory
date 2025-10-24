package service

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/service/model"
)

type IPartService interface {
	Get(ctx context.Context, uuid string) (model.Part, error)
	List(ctx context.Context, dto model.PartFilter) ([]model.Part, error)
	Create(ctx context.Context, dto model.Part) (string, error)
}
