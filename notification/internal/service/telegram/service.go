package telegram

import (
	"embed"
	"html/template"

	"github.com/Artyom099/factory/notification/internal/client/http"
	def "github.com/Artyom099/factory/notification/internal/service"
)

const chatID = 234586218

//go:embed templates/order_paid_notification.tmpl
var templateFS embed.FS

var orderPaidTemplate = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))

var _ def.ITelegramService = (*service)(nil)

type service struct {
	telegramClient http.ITelegramClient
}

func NewService(telegramClient http.ITelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}
