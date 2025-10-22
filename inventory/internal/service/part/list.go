package part

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *service) List(ctx context.Context, dto servModel.ModelPartFilter) ([]servModel.Part, error) {
	res, err := s.partRepository.List(ctx, converter.ModelToRepoPartFilter(dto))
	if err != nil {
		return []servModel.Part{}, err
	}

	return converter.RepoToModelListParts(res), nil
}
