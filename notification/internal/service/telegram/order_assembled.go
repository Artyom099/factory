package telegram

import (
	"bytes"
	"context"
	"text/template"

	"go.uber.org/zap"

	"github.com/Artyom099/factory/notification/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

var orderAssembleTemplate = template.Must(template.ParseFS(templateFS, "templates/order_assembled_notification.tmpl"))

func (s *service) SendOrderAssembledNotification(ctx context.Context, dto model.OrderAssembledInEvent) error {
	message, err := s.buildOrderAssembledMessage(dto)
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

func (s *service) buildOrderAssembledMessage(dto model.OrderAssembledInEvent) (string, error) {
	// todo

	var buf bytes.Buffer
	err := orderAssembleTemplate.Execute(&buf, dto)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
