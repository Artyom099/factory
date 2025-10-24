package part

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *service) Create(ctx context.Context, dto model.Part) (string, error) {
	partUUID, err := s.partRepository.Create(ctx, dto)
	if err != nil {
		return "", model.ErrInternalError
	}

	return partUUID, nil
}
