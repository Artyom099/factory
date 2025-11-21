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

func (r *repository) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	builderGet := sq.
		Select(
			"o.id AS order_uuid",
			"o.user_uuid",
			"COALESCE(array_agg(op.part_id), '{}') AS part_uuids",
			"o.total_price",
			"COALESCE(o.transaction_uuid, '') AS transaction_uuid",
			"o.payment_method",
			"o.status",
			"o.created_at",
			"o.updated_at",
		).
		From("orders o").
		LeftJoin("orders_parts op ON o.id = op.order_id").
		Where(sq.Eq{"o.id": orderUuid}).
		GroupBy(
			"o.id",
			"o.user_uuid",
			"o.total_price",
			"COALESCE(o.transaction_uuid, '')",
			"o.payment_method",
			"o.status",
			"o.created_at",
			"o.updated_at",
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderGet.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.Order{}, err
	}

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		log.Printf("failed to query order: %v\n", err)
		return model.Order{}, err
	}
	defer rows.Close()

	order, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[repoModel.RepoOrder])
	if err != nil {
		log.Printf("failed to scan order: %v\n", err)
		return model.Order{}, err
	}

	log.Printf("order: %v\n", order)

	return converter.ToModelOrder(order), nil
}
