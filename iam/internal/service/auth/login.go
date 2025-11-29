package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *service) Login(ctx context.Context, login, password string) (string, error) {
	// 1. Получаем пользователя из репозитория/бд
	user, err := s.userRepository.Get(ctx, login)
	if err != nil {
		return "", err
	}

	// 2. Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password)); err != nil {
		return "", model.ErrInvalidPassword
	}

	// 3. Генерируем sessionID
	sessionID := uuid.New().String()

	// 4. Создаём сессию (доменные структуры)
	session := model.Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
	}

	// 5. Сохраняем сессию в Redis
	err = s.sessionRepository.Create(ctx, session, *s.sessionTTL)
	if err != nil {
		return "", err
	}

	// 6. Добавляем sessionID в множество активных сессий пользователя
	err = s.sessionRepository.AddSessionToUserSet(ctx, user.ID, sessionID)
	if err != nil {
		return "", err
	}

	// 7. Возвращаем sessionID клиенту
	return sessionID, nil
}
