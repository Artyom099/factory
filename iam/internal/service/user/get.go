package user

import (
	"context"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Get(ctx context.Context, uuid string) (model.User, error) {
	user, err := s.userRepository.Get(ctx, uuid)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}
