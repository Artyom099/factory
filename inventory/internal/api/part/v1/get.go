package v1

import (
	"context"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return &inventoryV1.GetPartResponse{}, err
	}

	res, err := a.partService.Get(ctx, uuid.String())
	if err != nil {
		return &inventoryV1.GetPartResponse{}, err
	}

	return converter.ModelToApiPart(res), nil
}
