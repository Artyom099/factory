package v1

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	servResp, err := a.partService.Get(ctx, converter.PartGetApiRequestToPartGetServiceRequest(req))
	if err != nil {
		return nil, err
	}

	return converter.PartGetServiceResponseToPartGetApiResponse(servResp), nil
}
