package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Update(ctx context.Context, dto model.Order) error {
	order := converter.ToRepoOrder(dto)

	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("transaction_uuid", order.TransactionUUID).
		Set("payment_method", order.PaymentMethod).
		Set("status", order.Status).
		Where(sq.Eq{"order_uuid": order.OrderUUID})

	query, args, err := builderUpdate.ToSql()
	if err != nil {
		log.Fatalf("failed to build query: %v", err)
	}

	res, err := r.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v\n", err)
		return model.ErrInternalError
	}

	log.Printf("updated %d rows", res.RowsAffected())

	return nil
}
