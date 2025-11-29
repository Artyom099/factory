package v1

import (
	"github.com/Artyom099/factory/iam/internal/service"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

type api struct {
	userV1.UnimplementedUserServiceServer

	userServise service.IUserService
}

func NewAPI(userServise service.IUserService) *api {
	return &api{
		userServise: userServise,
	}
}
