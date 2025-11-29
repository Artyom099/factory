package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Create(ctx context.Context, dto model.Order) (string, error) {
	order := converter.ToRepoOrder(dto)

	builderInsertOrder := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("user_uuid", "total_price", "status", "payment_method").
		Values(order.UserUUID, order.TotalPrice, repoModel.OrderStatusPENDINGPAYMENT, repoModel.OrderPaymentMethodUNSPECIFIED).
		Suffix("RETURNING id")

	query, args, err := builderInsertOrder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return "", err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to insert order: %v\n", err)
		return "", err
	}
	defer rows.Close()

	orderUUID, err := pgx.CollectOneRow(rows, pgx.RowTo[string])
	if err != nil {
		log.Printf("failed to collect inserted order id: %v\n", err)
		return "", err
	}

	for _, partUuid := range order.PartUuids {
		builderInsertOrderParts := sq.Insert("orders_parts").
			Columns("order_id", "part_id").
			Values(orderUUID, partUuid).
			PlaceholderFormat(sq.Dollar)

		query, args, err := builderInsertOrderParts.ToSql()
		if err != nil {
			log.Printf("failed to build query: %v\n", err)
			return "", err
		}

		_, err = r.pool.Exec(ctx, query, args...)
		if err != nil {
			log.Printf("failed to insert order: %v\n", err)
			return "", err
		}
	}

	return orderUUID, nil
}
