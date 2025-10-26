package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Cancel(ctx context.Context, orderUUID string) error {
	builderUpdate := sq.Update("orders").
		PlaceholderFormat(sq.Dollar).
		Set("status", repoModel.OrderStatusCANCELLED).
		Where(sq.Eq{"order_uuid": orderUUID})

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
