package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"

	"github.com/Artyom099/factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger                 LoggerConfig
	OrderHTTP              OrderHTTPConfig
	Postgres               PostgresConfig
	Payment                PaymentClientConfig
	Inventory              InventoryClientConfig
	Kafka                  KafkaConfig
	OrderPaidProducer      OrderPaidProducerConfig
	OrderAssembledConsumer OrderAssembledConsumerConfig
	IamCLient              IamClientConfig
}

func Load(path ...string) error {
	dotenvDirs := make([]string, 0, len(path))
	for _, p := range path {
		if p == "" {
			continue
		}

		absPath := p
		if !filepath.IsAbs(p) {
			resolved, err := filepath.Abs(p)
			if err != nil {
				return err
			}
			absPath = resolved
		}

		dotenvDirs = append(dotenvDirs, filepath.Dir(absPath))
	}

	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := ensureAbsoluteMigrationsDir(dotenvDirs); err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
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

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderPaidProducerCfg, err := env.NewOrderPaidProducerConfig()
	if err != nil {
		return err
	}

	orderAssembledConsumerCfg, err := env.NewOrderAssembledConsumerConfig()
	if err != nil {
		return err
	}

	iamCLientCfg, err := env.NewIamClientConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:                 loggerCfg,
		OrderHTTP:              orderHTTPCfg,
		Postgres:               postgresCfg,
		Payment:                paymentCfg,
		Inventory:              inventoryCfg,
		Kafka:                  kafkaCfg,
		OrderPaidProducer:      orderPaidProducerCfg,
		OrderAssembledConsumer: orderAssembledConsumerCfg,
		IamCLient:              iamCLientCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}

func ensureAbsoluteMigrationsDir(dotenvDirs []string) error {
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" || filepath.IsAbs(migrationsDir) {
		return nil
	}

	if len(dotenvDirs) > 0 {
		absPath := filepath.Clean(filepath.Join(dotenvDirs[0], migrationsDir))
		serr := os.Setenv("MIGRATIONS_DIR", absPath)
		if serr != nil {
			return serr
		}

		return nil
	}

	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return err
	}

	err = os.Setenv("MIGRATIONS_DIR", absPath)
	if err != nil {
		return err
	}

	return nil
}
