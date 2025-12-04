package env

import (
	"net"

	"github.com/caarlos0/env/v11"
)

type iamCLientEnvConfig struct {
	Host string `env:"IAM_CLIENT_HOST,required"`
	Port string `env:"IAM_CLIENT_PORT,required"`
}

type iamCLientConfig struct {
	raw iamCLientEnvConfig
}

func NewIamClientConfig() (*iamCLientConfig, error) {
	var raw iamCLientEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &iamCLientConfig{raw: raw}, nil
}

func (cfg *iamCLientConfig) Address() string {
	return net.JoinHostPort(cfg.raw.Host, cfg.raw.Port)
}
