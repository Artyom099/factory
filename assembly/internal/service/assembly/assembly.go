package assembly

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"go.uber.org/zap"

	"github.com/Artyom099/factory/assembly/internal/model"
	"github.com/Artyom099/factory/platform/pkg/logger"
)

func (s *service) Assembly(ctx context.Context, dto model.OrderPaidInEvent) error {
	// Генерируем безопасную случайную задержку 5–10 сек
	n, err := rand.Int(rand.Reader, big.NewInt(6))
	if err != nil {
		return fmt.Errorf("failed to generate delay: %w", err)
	}
	buildTime := time.Duration(5+n.Int64()) * time.Second
	logger.Debug(ctx, "Build Time", zap.Any("buildTime - ", buildTime))

	// Ждём buildTime, но с учётом отмены контекста
	timer := time.NewTimer(buildTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		// Заказ собран
	case <-ctx.Done():
		return ctx.Err()
	}

	logger.Debug(ctx, "OrderPaidInEvent", zap.Any("dto - ", dto))

	// отправляем сообщение в кафку, что заказ собран
	err = s.assemblyProducerService.ProduceOrderAssembled(ctx, model.OrderAssembledOutEvent{
		EventUUID:    dto.EventUUID,
		OrderUUID:    dto.OrderUUID,
		UserUUID:     dto.UserUUID,
		BuildTimeSec: int64(buildTime.Seconds()),
	})
	if err != nil {
		return model.ErrSendOderAssembledMessageToKafka
	}

	return nil
}
