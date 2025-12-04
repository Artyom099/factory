package session

import (
	"fmt"

	def "github.com/Artyom099/factory/iam/internal/repository"
	"github.com/Artyom099/factory/platform/pkg/cache"
)

const (
	cacheKeyPrefix = "iam:session:"
)

var _ def.ISessionRepository = (*repository)(nil)

type repository struct {
	cache cache.RedisClient
}

func NewRepository(cache cache.RedisClient) *repository {
	return &repository{
		cache: cache,
	}
}

// api -> service -> cache repo -> redis client (обертка наша) -> redis
func (r *repository) getCacheKey(uuid string) string {
	return fmt.Sprintf("%s%s", cacheKeyPrefix, uuid)
}

func (r *repository) getUserSessionsKey(userID string) string {
	return fmt.Sprintf("user:%s:sessions", userID)
}
