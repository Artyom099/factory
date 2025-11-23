package telegram

import (
	"bytes"
	"context"

	"go.uber.org/zap"

	"github.com/Artyom099/factory/notification/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

func (s *service) SendOrderPaidNotification(ctx context.Context, dto model.OrderPaidInEvent) error {
	message, err := s.buildOrderPaidMessage(dto)
	if err != nil {
		return err
	}

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
