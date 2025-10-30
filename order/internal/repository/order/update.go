package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Update(ctx context.Context, dto model.Order) error {
	order := converter.ToRepoOrder(dto)

	builderUpdate := sq.Update("orders").
		Set("transaction_uuid", order.TransactionUUID).
		Set("payment_method", order.PaymentMethod).
		Set("status", order.Status).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": order.OrderUUID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v", err)
		return err
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v\n", err)
		return err
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil
}
