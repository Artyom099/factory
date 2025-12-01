package session

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/iam/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, sessionUUID string) (model.Session, error) {
	cacheKey := r.getCacheKey(sessionUUID)

	values, err := r.cache.HGetAll(ctx, cacheKey)
	if err != nil {
		if errors.Is(err, redigo.ErrNil) {
			return model.Session{}, model.ErrSessionNotFound
		}
		return model.Session{}, err
	}

	if len(values) == 0 {
		return model.Session{}, model.ErrSessionNotFound
	}

	var sessionRedisView repoModel.SessionRedisView
	err = redigo.ScanStruct(values, &sessionRedisView)
	if err != nil {
		return model.Session{}, err
	}

	return converter.ToModelSession(sessionRedisView), nil
}
