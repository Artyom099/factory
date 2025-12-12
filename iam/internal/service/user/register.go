package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Register(ctx context.Context, user model.User) (string, error) {
	ctx, getUserSpan := tracing.StartSpan(ctx, "iam.get_user",
		trace.WithAttributes(
			attribute.String("user.login", user.Login),
		),
	)

	// Проверяем, что пользователь с таким email не существует
	_, err := s.userRepository.Get(ctx, user.Email)

	// Если вернулся пользователь — он уже существует
	if err == nil {
		getUserSpan.End()
		return "", model.ErrUserAlreadyExists
	}
	// Если это НЕ ошибка "не найден" — вернуть ошибку
	if !errors.Is(err, model.ErrUserNotFound) {
		getUserSpan.RecordError(err)
		getUserSpan.End()
		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	getUserSpan.End()

	ctx, createUserSpan := tracing.StartSpan(ctx, "iam.create_user",
		trace.WithAttributes(
			attribute.String("user.login", user.Login),
		),
	)
	defer createUserSpan.End()

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		createUserSpan.RecordError(err)
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Hash = string(hashed)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = &user.CreatedAt

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		createUserSpan.RecordError(err)
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
