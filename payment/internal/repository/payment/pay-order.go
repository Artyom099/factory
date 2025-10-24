package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/payment/internal/service/model"
)

func (r *repository) PayOrder(ctx context.Context, dto model.Payment) (string, error) {
	transactionUuid := uuid.New().String()

	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUuid)

	return transactionUuid, nil
}
