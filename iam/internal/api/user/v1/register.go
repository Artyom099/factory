package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/iam/internal/api/converter"
	"github.com/Artyom099/factory/platform/pkg/logger"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

func (a *api) Register(ctx context.Context, req *userV1.RegisterRequest) (*userV1.RegisterResponse, error) {
	err := req.Validate()
	if err != nil {
		logger.Error(ctx, "validation error", zap.String("login", req.GetInfo().GetInfo().Login))
		return &userV1.RegisterResponse{}, status.Errorf(codes.InvalidArgument, "validation error")
	}

	err = a.userServise.Register(ctx, converter.ToModelRegUser(req))
	if err != nil {
		logger.Error(ctx, "internal server error", zap.Error(err))
		return &userV1.RegisterResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return nil, nil
}
