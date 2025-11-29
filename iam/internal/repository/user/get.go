package user

import (
	"context"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (r *repository) Get(ctx context.Context, login string) (model.User, error) {
	return model.User{}, nil
}
