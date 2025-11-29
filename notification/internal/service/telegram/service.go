package telegram

import (
	"embed"

	"github.com/Artyom099/factory/notification/internal/client/http"
	def "github.com/Artyom099/factory/notification/internal/service"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

var _ def.INotificationTelegramService = (*service)(nil)

type service struct {
	telegramClient http.ITelegramClient
	chatID         int64
}

func NewService(telegramClient http.ITelegramClient, chatID int64) *service {
	return &service{
		telegramClient: telegramClient,
		chatID:         chatID,
	}
}
