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
	IamCLient   IamClientConfig
	Tracing     TracingConfig
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

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	iamCLientCfg, err := env.NewIamClientConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		PaymentGRPC: paymentGRPCCfg,
		IamCLient:   iamCLientCfg,
		Tracing:     tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
