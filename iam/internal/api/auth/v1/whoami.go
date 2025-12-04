package v1

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/iam/internal/api/converter"
	"github.com/Artyom099/factory/platform/pkg/logger"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
)

func (a *api) Whoami(ctx context.Context, req *authV1.WhoamiRequest) (*authV1.WhoamiResponse, error) {
	err := req.Validate()
	if err != nil {
		logger.Error(ctx, "validation error", zap.String("session_uuid", req.GetSessionUuid()))
		return &authV1.WhoamiResponse{}, status.Errorf(codes.InvalidArgument, "validation error")
	}

	user, session, err := a.authServise.Whoami(ctx, req.GetSessionUuid())
	if err != nil {
		logger.Error(ctx, "internal server error", zap.Error(err))
		return &authV1.WhoamiResponse{}, status.Errorf(codes.Internal, "internal server error")
	}

	return converter.ToApiWhoami(user, session), nil
}
