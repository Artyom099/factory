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

	// kafkaCfg, err := env.NewKafkaConfig()
	// if err != nil {
	// 	return err
	// }

	// orderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	// if err != nil {
	// 	return err
	// }

	// orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	// if err != nil {
	// 	return err
	// }

	appConfig = &config{
		IamGRPC:  iamGRPCCfg,
		Logger:   loggerCfg,
		Postgres: postgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
