package part

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *service) List(ctx context.Context, dto model.PartFilter) ([]model.Part, error) {
	res, err := s.partRepository.List(ctx, dto)
	if err != nil {
		return []model.Part{}, err
	}

	return res, nil
}
