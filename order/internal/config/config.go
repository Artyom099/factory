package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Artyom099/factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger    LoggerConfig
	OrderGRPC OrderGRPCConfig
	Postgres  PostgresConfig
	Payment   PaymentClientConfig
	Inventory InventoryClientConfig
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

	orderGRPCCfg, err := env.NewOrderGRPCConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	paymentCfg, err := env.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	inventoryCfg, err := env.NewInventoryClientConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:    loggerCfg,
		OrderGRPC: orderGRPCCfg,
		Postgres:  postgresCfg,
		Payment:   paymentCfg,
		Inventory: inventoryCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
