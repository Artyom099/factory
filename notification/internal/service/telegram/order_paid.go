package telegram

import (
	"bytes"
	"context"
	"text/template"

	"go.uber.org/zap"

	"github.com/Artyom099/factory/notification/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

var orderPaidTemplate = template.Must(template.ParseFS(templateFS, "templates/order_paid_notification.tmpl"))

func (s *service) SendOrderPaidNotification(ctx context.Context, dto model.OrderPaidInEvent) error {
	message, err := s.buildOrderPaidMessage(dto)
	if err != nil {
		return err
	}

	err = s.telegramClient.SendMessage(ctx, s.chatID, message)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Telegram message sent to chat", zap.Int64("chat_id", s.chatID), zap.String("message", message))
	return nil
}

func (s *service) buildOrderPaidMessage(dto model.OrderPaidInEvent) (string, error) {
	var buf bytes.Buffer
	err := orderPaidTemplate.Execute(&buf, dto)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
