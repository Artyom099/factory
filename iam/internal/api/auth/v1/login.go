package v1

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/iam/internal/api/converter"
	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
)

func (a *api) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	err := req.Validate()
	if err != nil {
		logger.Error(ctx, "validation error", zap.String("login", req.GetLogin()))
		return &authV1.LoginResponse{}, status.Errorf(codes.InvalidArgument, "validation error")
	}

	sessionUUID, err := a.authServise.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		if errors.Is(err, model.ErrInvalidPassword) {
			logger.Error(ctx, "wrong password", zap.String("login", req.GetLogin()))
			return &authV1.LoginResponse{}, status.Errorf(codes.Unauthenticated, "wrong password")
		}

		logger.Error(ctx, "internal server error", zap.Error(err))
		return &authV1.LoginResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return converter.ToApiLogin(sessionUUID), nil
}
