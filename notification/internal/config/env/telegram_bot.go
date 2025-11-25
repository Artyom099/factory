package env

import "github.com/caarlos0/env/v11"

type telegramEnvConfig struct {
	Token  string `env:"TELEGRAM_BOT_TOKEN,required"`
	ChatID int64  `env:"TELEGRAM_BOT_CHAT_ID,required"`
}

type telegramConfig struct {
	raw telegramEnvConfig
}

func NewTelegramConfig() (*telegramConfig, error) {
	var raw telegramEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &telegramConfig{raw: raw}, nil
}

func (cfg *telegramConfig) Token() string {
	return cfg.raw.Token
}

func (cfg *telegramConfig) ChatID() int64 {
	return cfg.raw.ChatID
}
