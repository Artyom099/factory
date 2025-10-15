package part

import (
	"context"

	servModel "github.com/Artyom099/factory/inventory/internal/model"
	"github.com/Artyom099/factory/inventory/internal/repository/converter"
)

func (s *service) List(ctx context.Context, dto servModel.PartListServiceRequest) (servModel.PartListServiceResponse, error) {
	res, err := s.partRepository.List(ctx, converter.PartListServiceRequestToPartListRepoRequest(dto))
	if err != nil {
		return servModel.PartListServiceResponse{}, err
	}

	return converter.PartListRepoResponseToPartListServiceResponse(res), nil
}
