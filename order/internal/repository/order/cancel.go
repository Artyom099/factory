package order

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
)

func (r *repository) Cancel(ctx context.Context, orderUUID string) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", repoModel.OrderStatusCANCELLED).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": orderUUID})

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
