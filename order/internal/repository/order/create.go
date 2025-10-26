package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Create(ctx context.Context, dto model.Order) (string, error) {
	order := converter.ToRepoOrder(dto)

	builderInsert := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "part_uuids", "total_price", "status", "payment_method").
		Values(order.UserUUID, order.PartUuids, order.TotalPrice, repoModel.OrderStatusPENDINGPAYMENT, repoModel.OrderPaymentMethodUNSPECIFIED).
		Suffix("RETURNING order_uuid")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return "", model.ErrInternalError
	}

	var orderUUID string
	err = r.pool.QueryRow(ctx, query, args...).Scan(&orderUUID)
	if err != nil {
		log.Printf("failed to insert order: %v\n", err)
		return "", model.ErrInternalError
	}

	return orderUUID, nil
}
