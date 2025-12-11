package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/iam/internal/api/converter"
	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

func (a *api) Get(ctx context.Context, req *userV1.GetUserRequest) (*userV1.GetUserResponse, error) {
	userUUID, err := uuid.Parse(req.GetUserUuid())
	if err != nil {
		logger.Error(ctx, "invalid user_uuid", zap.String("user_uuid", req.GetUserUuid()))
		return &userV1.GetUserResponse{}, status.Errorf(codes.InvalidArgument, "invalid uuid")
	}

	user, err := a.userServise.Get(ctx, userUUID.String())
	if err != nil {
		if errors.Is(err, model.ErrUserNotFound) {
			logger.Error(ctx, "user not found", zap.String("user_uuid", userUUID.String()))
			return &userV1.GetUserResponse{}, status.Errorf(codes.NotFound, "user not found")
		}

		logger.Error(ctx, "internal server error", zap.Error(err))
		return &userV1.GetUserResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return converter.ToApiUser(user), nil
}
