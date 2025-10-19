package v1

import (
	"context"

	clientConverter "github.com/Artyom099/factory/order/internal/client/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.ListPartsFilter) (model.ListPartsResponseDto, error) {
	parts, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.PartFilterToProto(filter),
	})
	if err != nil {
		return model.ListPartsResponseDto{}, err
	}

	return clientConverter.PartFilterToModel(parts.Parts), nil
}
