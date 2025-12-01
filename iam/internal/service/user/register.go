package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Register(ctx context.Context, user model.User) (string, error) {
	// Проверяем, что пользователь с таким email не существует
	_, err := s.userRepository.Get(ctx, user.Email)
	if err != nil && !errors.Is(err, model.ErrUserNotFound) {
		if errors.Is(err, model.ErrUserNotFound) {
			return "", model.ErrUserNotFound
		}

		return "", fmt.Errorf("failed to check existing user: %w", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}

	user.Hash = string(hashed)
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = &user.CreatedAt

	id, err := s.userRepository.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	return id, nil
}
