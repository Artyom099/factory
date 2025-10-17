package v1

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	res, err := a.partService.List(ctx, converter.PartListApiRequestToPartGetServiceRequest(req))
	if err != nil {
		return &inventoryV1.ListPartsResponse{}, err
	}

	return converter.PartListServiceResponseToPartGetApiResponse(res), nil
}
