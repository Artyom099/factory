package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/inventory/internal/api/converter"
	"github.com/Artyom099/factory/platform/pkg/logger"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	res, err := a.partService.List(ctx, converter.ToModelPartFilter(req))
	if err != nil {
		logger.Error(ctx, "internal server error", zap.Error(err))
		return &inventoryV1.ListPartsResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return converter.ToApiListParts(res), nil
}
