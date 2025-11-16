package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type asemblyGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type asemblyGRPCConfig struct {
	raw asemblyGRPCEnvConfig
}

func NewAsemblyGRPCConfig() (*asemblyGRPCConfig, error) {
	var raw asemblyGRPCEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &asemblyGRPCConfig{raw: raw}, nil
}

func (cfg *asemblyGRPCConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
