package auth

import (
	"context"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Whoami(ctx context.Context, sessionUuid string) (model.User, model.Session, error) {
	session, err := s.sessionRepository.Get(ctx, sessionUuid)
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	user, err := s.userRepository.Get(ctx, session.UserID)
	if err != nil {
		return model.User{}, model.Session{}, err
	}

	return user, session, nil
}
