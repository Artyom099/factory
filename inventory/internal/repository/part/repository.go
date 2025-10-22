package part

import (
	"sync"

	def "github.com/Artyom099/factory/inventory/internal/repository"
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
)

var _ def.IPartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.RepoPart
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.RepoPart),
	}
}
