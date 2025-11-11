package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Artyom099/factory/payment/internal/config/env"
)

var appConfig *config

type config struct {
	Logger      LoggerConfig
	PaymentGRPC PaymentGRPCConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewUFOGRPCConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: inventoryGRPCCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
