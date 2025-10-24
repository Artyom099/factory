package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	"github.com/Artyom099/factory/inventory/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		return &inventoryV1.GetPartResponse{}, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}

	res, err := a.partService.Get(ctx, uuid.String())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &inventoryV1.GetPartResponse{}, status.Errorf(codes.NotFound, "order not found")
		}
		return &inventoryV1.GetPartResponse{}, status.Errorf(codes.Internal, "intermal server error")
	}

	return converter.ToApiPart(res), nil
}
