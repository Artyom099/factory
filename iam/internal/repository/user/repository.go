package user

import (
	"github.com/jackc/pgx/v5/pgxpool"

	def "github.com/Artyom099/factory/iam/internal/repository"
)

var _ def.IUserRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pool: pool,
	}
}
