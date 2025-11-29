package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type IamGRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	DatabaseName() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}

type SessionConfig interface {
	TTL() string
}
