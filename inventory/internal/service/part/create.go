package part

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *service) Create(ctx context.Context, dto servModel.Part) (string, error) {
	partUUID, err := s.partRepository.Create(ctx, converter.ModelToRepoPart(dto))
	if err != nil {
		return "", err
	}

	return partUUID, nil
}
