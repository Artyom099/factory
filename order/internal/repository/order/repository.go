package order

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/Artyom099/factory/order/internal/repository"
)

var _ def.IOrderRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
