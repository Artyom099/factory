package part

import (
	"context"
	"errors"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

func (s *service) Get(ctx context.Context, uuid string) (servModel.Part, error) {
	res, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		if errors.Is(err, repoModel.ErrPartNotFound) {
			return servModel.Part{}, servModel.ErrPartNotFound
		}
		return servModel.Part{}, err
	}

	return res, nil
}
