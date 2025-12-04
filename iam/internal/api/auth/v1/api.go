package v1

import (
	"github.com/Artyom099/factory/iam/internal/service"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
)

type api struct {
	authV1.UnimplementedAuthServiceServer

	authServise service.IAuthService
}

func NewAPI(partService service.IAuthService) *api {
	return &api{
		authServise: partService,
	}
}
