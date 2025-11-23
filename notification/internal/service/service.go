package service

import (
	"context"

	"github.com/Artyom099/factory/notification/internal/model"
)

type INotificationConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ITelegramService interface {
	SendOrderPaidNotification(ctx context.Context, dto model.OrderPaidInEvent) error
	SendOrderAssembledNotification(ctx context.Context, dto model.OrderAssembledInEvent) error
}
