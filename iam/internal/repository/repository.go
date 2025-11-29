package repository

import (
	"context"
	"time"

	"github.com/Artyom099/factory/iam/internal/model"
)

type ISessionRepository interface {
	Get(ctx context.Context, sessionUUID string) (model.Session, error)
	Create(ctx context.Context, session model.Session, ttl time.Duration) error
	AddSessionToUserSet(ctx context.Context, userID, sessionID string) error
}

type IUserRepository interface {
	Get(ctx context.Context, login string) (model.User, error)
	Create(ctx context.Context) error
}
