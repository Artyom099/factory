package auth

import (
	"context"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUuid string) (model.User, error) {
	return model.User{}, nil
}
