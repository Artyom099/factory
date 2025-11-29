package app

import (
	"context"
	"log"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"

	authApiV1 "github.com/Artyom099/factory/iam/internal/api/auth/v1"
	userV1API "github.com/Artyom099/factory/iam/internal/api/user/v1"
	"github.com/Artyom099/factory/iam/internal/config"
	"github.com/Artyom099/factory/iam/internal/repository"
	sessionRepository "github.com/Artyom099/factory/iam/internal/repository/session"
	userRepository "github.com/Artyom099/factory/iam/internal/repository/user"
	"github.com/Artyom099/factory/iam/internal/service"
	authService "github.com/Artyom099/factory/iam/internal/service/auth"
	userService "github.com/Artyom099/factory/iam/internal/service/user"
	"github.com/Artyom099/factory/platform/pkg/cache"
	"github.com/Artyom099/factory/platform/pkg/cache/redis"
	"github.com/Artyom099/factory/platform/pkg/logger"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

type diContainer struct {
	authV1API         authV1.AuthServiceServer
	authService       service.IAuthService
	sessionRepository repository.ISessionRepository

	userV1API      userV1.UserServiceServer
	userService    service.IUserService
	userRepository repository.IUserRepository

	sessionTTL *time.Duration

	postgresHandle *pgxpool.Pool

	redisPool   *redigo.Pool
	redisClient cache.RedisClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// auth

func (d *diContainer) AuthV1API(ctx context.Context) authV1.AuthServiceServer {
	if d.authV1API == nil {
		d.authV1API = authApiV1.NewAPI(d.AuthService(ctx))
	}

	return d.authV1API
}

func (d *diContainer) AuthService(ctx context.Context) service.IAuthService {
	if d.authService == nil {
		d.authService = authService.NewService(
			d.Session(),
			d.UserRepository(ctx),
			d.SessionRepository(),
		)
	}

	return d.authService
}

// session

func (d *diContainer) SessionRepository() repository.ISessionRepository {
	if d.sessionRepository == nil {
		d.sessionRepository = sessionRepository.NewRepository(d.RedisClient())
	}

	return d.sessionRepository
}

// user

func (d *diContainer) UserV1API(ctx context.Context) userV1.UserServiceServer {
	if d.userV1API == nil {
		d.userV1API = userV1API.NewAPI(d.UserService(ctx))
	}

	return d.userV1API
}

func (d *diContainer) UserService(ctx context.Context) service.IUserService {
	if d.userService == nil {
		d.userService = userService.NewService(d.UserRepository(ctx))
	}

	return d.userService
}

func (d *diContainer) UserRepository(ctx context.Context) repository.IUserRepository {
	if d.userRepository == nil {
		d.userRepository = userRepository.NewRepository(d.PostgresHandle(ctx))
	}

	return d.userRepository
}

// pg & redis

func (d *diContainer) PostgresHandle(ctx context.Context) *pgxpool.Pool {
	if d.postgresHandle == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Fatalf("failed to connect postgres db: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping postgres db: %v\n", err)
		}

		d.postgresHandle = pool
	}

	return d.postgresHandle
}

func (d *diContainer) RedisPool() *redigo.Pool {
	if d.redisPool == nil {
		d.redisPool = &redigo.Pool{
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
		}
	}

	return d.redisPool
}

func (d *diContainer) RedisClient() cache.RedisClient {
	if d.redisClient == nil {
		d.redisClient = redis.NewClient(d.RedisPool(), logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
	}

	return d.redisClient
}

func (d *diContainer) Session() *time.Duration {
	if d.sessionTTL == nil {
		ttlStr := config.AppConfig().Session.TTL()

		parsed, err := time.ParseDuration(ttlStr)
		if err != nil {
			log.Fatalf("failed to parse ttl: %v", err)
		}

		d.sessionTTL = &parsed
	}

	return d.sessionTTL
}
