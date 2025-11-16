package telegram

import (
	"bytes"
	"context"
	"embed"
	"html/template"

	"github.com/Artyom099/factory/notification/internal/client/http"
	"github.com/Artyom099/factory/notification/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
	"go.uber.org/zap"
)

//go:embed templates/order_paid_notification.tmpl
var templateFS embed.FS

var orderPaidTemplate = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))

type service struct {
	telegramClient http.TelegramClient
}

func NewService(telegramClient http.TelegramClient) *service {
	return &service{
		telegramClient: telegramClient,
	}
}

func (s *service) SendOrderPaidNotification(ctx context.Context, dto model.OrderPaidInEvent) error {
	message, err := s.buildOrderPaidMessage(dto)
	if err != nil {
		return err
	}

	var chatID int64 = 123

	err = s.telegramClient.SendMessage(ctx, chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(dto model.OrderPaidInEvent) (string, error) {
	// todo

	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, dto)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
