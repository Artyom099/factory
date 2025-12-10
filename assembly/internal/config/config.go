package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Artyom099/factory/assembly/internal/config/env"
)

var appConfig *config

type config struct {
	AssemblyGRPC           AssemblyGRPCConfig
	Logger                 LoggerConfig
	Kafka                  KafkaConfig
	OrderAssembledProducer OrderAssembledProducerConfig
	OrderPaidConsumer      OrderPaidConsumerConfig
	Tracing                TracingConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	assemblyGRPCCfg, err := env.NewAsemblyGRPCConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	orderAssembledProducerCfg, err := env.NewOrderAssembledProducerConfig()
	if err != nil {
		return err
	}

	orderPaidConsumerCfg, err := env.NewOrderPaidConsumerConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		AssemblyGRPC:           assemblyGRPCCfg,
		Logger:                 loggerCfg,
		Kafka:                  kafkaCfg,
		OrderAssembledProducer: orderAssembledProducerCfg,
		OrderPaidConsumer:      orderPaidConsumerCfg,
		Tracing:                tracingCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
