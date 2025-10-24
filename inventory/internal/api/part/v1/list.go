package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	res, err := a.partService.List(ctx, converter.ToModelPartFilter(req))
	if err != nil {
		return &inventoryV1.ListPartsResponse{}, status.Errorf(codes.Internal, "intermal server error")
	}

	return converter.ToApiListParts(res), nil
}
