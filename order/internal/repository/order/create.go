package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/order/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, dto model.OrderCreateRepoRequestDto) (string, error) {
	orderUuid := uuid.New().String()

	order := model.OrderGetRepoResponseDto{
		OrderUUID:       orderUuid,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: "",
		PaymentMethod:   "",
		Status:          model.OrderStatus("PENDING_PAYMENT"),
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[orderUuid] = &order

	return orderUuid, nil
}
