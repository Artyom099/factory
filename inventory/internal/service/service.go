package service

import (
	"context"

	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

type IPartService interface {
	Get(ctx context.Context, uuid string) (servModel.Part, error)
	List(ctx context.Context, dto servModel.ModelPartFilter) ([]servModel.Part, error)
	Create(ctx context.Context, dto servModel.Part) (string, error)
}
