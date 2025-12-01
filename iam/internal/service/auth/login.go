package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	user, err := s.userRepository.Get(ctx, login)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return "", model.ErrInvalidPassword
	}

	sessionID := uuid.New().String()

	session := model.Session{
		ID:        sessionID,
		UserID:    user.ID,
		Login:     user.Login,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(*s.sessionTTL),
	}

	// Сохраняем сессию в Redis
	err = s.sessionRepository.Create(ctx, session, *s.sessionTTL)
	if err != nil {
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.ID, sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}
