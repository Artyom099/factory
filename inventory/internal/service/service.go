package service

import (
	"context"

	servModel "github.com/Artyom099/factory/inventory/internal/model"
)

type IPartService interface {
	Get(ctx context.Context, dto servModel.PartGetServiceRequest) (servModel.PartGetServiceResponse, error)
	List(ctx context.Context, dto servModel.PartListServiceRequest) (servModel.PartListServiceResponse, error)
	Create(ctx context.Context, dto servModel.PartCreateServiceRequest) (string, error)
}
