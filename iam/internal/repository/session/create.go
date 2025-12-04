package session

import (
	"context"
	"time"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/iam/internal/repository/converter"
)

func (r *repository) Create(ctx context.Context, session model.Session, ttl time.Duration) error {
	cacheKey := r.getCacheKey(session.ID)

	redisView := converter.ToRedisViewSession(session)

	err := r.cache.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		return err
	}

	return r.cache.Expire(ctx, cacheKey, ttl)
}
