package part

import (
	"context"

	servModel "github.com/Artyom099/factory/inventory/internal/model"
	"github.com/Artyom099/factory/inventory/internal/repository/converter"
)

func (s *service) Create(ctx context.Context, dto servModel.PartCreateServiceRequest) (string, error) {
	uuid, err := s.partRepository.Create(ctx, converter.PartCreateServiceRequestToPartCreateRepoRequest(dto))
	if err != nil {
		return "", err
	}

	return uuid, nil
}
