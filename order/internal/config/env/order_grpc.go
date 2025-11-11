package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type orderGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type orderGRPCConfig struct {
	raw orderGRPCEnvConfig
}

func NewOrderGRPCConfig() (*orderGRPCConfig, error) {
	var raw orderGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderGRPCConfig{raw: raw}, nil
}

func (cfg *orderGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
