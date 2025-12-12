package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	ctx, getUserSpan := tracing.StartSpan(ctx, "iam.get_user",
		trace.WithAttributes(
			attribute.String("user.login", login),
		),
	)

	user, err := s.userRepository.Get(ctx, login)
	if err != nil {
		getUserSpan.RecordError(err)
		getUserSpan.End()
		return "", err
	}

	getUserSpan.End()

	ctx, createSessionSpan := tracing.StartSpan(ctx, "iam.create_session",
		trace.WithAttributes(
			attribute.String("user.login", login),
		),
	)
	defer createSessionSpan.End()

	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		createSessionSpan.RecordError(err)
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
		createSessionSpan.RecordError(err)
		return "", err
	}

	err = s.sessionRepository.AddSessionToUserSet(ctx, user.ID, sessionID)
	if err != nil {
		createSessionSpan.RecordError(err)
		return "", err
	}

	return sessionID, nil
}
