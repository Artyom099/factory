package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Artyom099/factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	IamGRPC  IamGRPCConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	Session  SessionConfig
	Tracing  TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	iamGRPCCfg, err := env.NewIamGRPCConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	sessionCfg, err := env.NewSessionConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		IamGRPC:  iamGRPCCfg,
		Logger:   loggerCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
		Session:  sessionCfg,
		Tracing:  tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
