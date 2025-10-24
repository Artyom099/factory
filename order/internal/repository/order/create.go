package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Create(ctx context.Context, dto model.Order) (string, error) {
	orderUuid := uuid.New().String()

	order := converter.ToRepoOrder(dto)
	order.OrderUUID = orderUuid
	order.Status = repoModel.OrderStatusPENDINGPAYMENT

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[orderUuid] = &order

	return orderUuid, nil
}
