package part

import (
	"github.com/Artyom099/factory/inventory/internal/repository"
	def "github.com/Artyom099/factory/inventory/internal/service"
)

var _ def.IPartService = (*service)(nil)

type service struct {
	partRepository repository.IPartRepository
}

func NewService(partRepository repository.IPartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
