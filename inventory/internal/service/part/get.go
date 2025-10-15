package part

import (
	"context"

	servModel "github.com/Artyom099/factory/inventory/internal/model"
	"github.com/Artyom099/factory/inventory/internal/repository/converter"
)

func (s *service) Get(ctx context.Context, dto servModel.PartGetServiceRequest) (servModel.PartGetServiceResponse, error) {
	res, err := s.partRepository.Get(ctx, converter.PartGetServiceRequestToPartGetRepoRequest(dto))
	if err != nil {
		return servModel.PartGetServiceResponse{}, err
	}

	return converter.PartGetRepoResponseToPartGetServiceResponse(res), nil
}
