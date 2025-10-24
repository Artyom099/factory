package order

import (
	"sync"

	def "github.com/Artyom099/factory/order/internal/repository"
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
)

var _ def.IOrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]*repoModel.RepoOrder
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]*repoModel.RepoOrder),
	}
}
