package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	"github.com/Artyom099/factory/inventory/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	uuid, err := uuid.Parse(req.GetUuid())
	if err != nil {
		logger.Error(ctx, "invalid part_uuid", zap.String("part_uuid", req.GetUuid()))
		return &inventoryV1.GetPartResponse{}, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}

	res, err := a.partService.Get(ctx, uuid.String())
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			logger.Error(ctx, "part not found", zap.String("part_uuid", uuid.String()))
			return &inventoryV1.GetPartResponse{}, status.Errorf(codes.NotFound, "part not found")
		}

		logger.Error(ctx, "internal server error", zap.String("part_uuid", uuid.String()), zap.Error(err))
		return &inventoryV1.GetPartResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return converter.ToApiPart(res), nil
}
