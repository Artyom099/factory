package assembly

import (
	def "github.com/Artyom099/factory/assembly/internal/service"
)

var _ def.IAssemblyService = (*service)(nil)

type service struct {
	assemblyProducerService def.IAssemblyProducerService
}

func NewService(
	assemblyProducerService def.IAssemblyProducerService,
) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
	}
}
