package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/Artyom099/factory/order/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (r *repository) Get(ctx context.Context, orderUuid string) (model.Order, error) {
	builderGet := sq.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		From("orders").
		Where(sq.Eq{"order_uuid": orderUuid}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builderGet.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return model.Order{}, model.ErrInternalError
	}

	var order repoModel.RepoOrder

	err = r.pool.QueryRow(ctx, query, args...).
		Scan(&order.OrderUUID, &order.UserUUID, &order.PartUuids, &order.TotalPrice, &order.TransactionUUID, &order.PaymentMethod, &order.Status)
	if err != nil {
		log.Printf("failed to scan order: %v\n", err)
		return model.Order{}, model.ErrInternalError
	}

	log.Printf("order: %v\n", order)

	return converter.ToModelOrder(order), nil
}
