package v1

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) CreatePart(ctx context.Context, req *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	uuid, err := a.partService.Create(ctx, converter.PartCreateApiRequestToPartGetServiceRequest(req))
	if err != nil {
		return nil, err
	}

	return &inventoryV1.CreatePartResponse{Uuid: uuid}, nil
}
