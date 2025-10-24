package v1

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) CreatePart(ctx context.Context, req *inventoryV1.CreatePartRequest) (*inventoryV1.CreatePartResponse, error) {
	uuid, err := a.partService.Create(ctx, converter.ToModelPart(req))
	if err != nil {
		return &inventoryV1.CreatePartResponse{}, status.Errorf(codes.Internal, "intermal server error")
	}

	return &inventoryV1.CreatePartResponse{Uuid: uuid}, nil
}
