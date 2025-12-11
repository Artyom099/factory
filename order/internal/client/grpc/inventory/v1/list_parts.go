package v1

import (
	"context"

	clientConverter "github.com/Artyom099/factory/order/internal/client/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
	grpcAuth "github.com/Artyom099/factory/platform/pkg/middleware/grpc"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.ListPartsFilter) ([]model.Part, error) {
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	parts, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.PartFilterToProto(filter),
	})
	if err != nil {
		return []model.Part{}, err
	}

	return clientConverter.PartFilterToModel(parts.Parts), nil
}
