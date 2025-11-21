package service

import (
	"context"

	"github.com/Artyom099/factory/notification/internal/model"
)

type IConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, dto model.OrderPaidInEvent) error
}
