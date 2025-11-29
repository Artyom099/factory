package service

import (
	"context"

	"github.com/Artyom099/factory/iam/internal/model"
)

type IAuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Whoami(ctx context.Context, sessionUuid string) (model.User, error)
}

type IUserService interface {
	Get(ctx context.Context, uuid string) error
	Register(ctx context.Context) error
}
