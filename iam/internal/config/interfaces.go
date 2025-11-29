package config

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

type RedisConfig interface{}

type SessionConfig interface{}
