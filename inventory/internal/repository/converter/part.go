package converter

import (
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

// converter service - repo

func PartGetServiceRequestToPartGetRepoRequest(req servModel.PartGetServiceRequest) repoModel.PartGetRepoRequest {
	return repoModel.PartGetRepoRequest{Uuid: req.Uuid}
}

func PartGetRepoResponseToPartGetServiceResponse(dto repoModel.PartGetRepoResponse) servModel.PartGetServiceResponse {
	var dims *servModel.Dimensions
	if dto.Part.Dimensions != nil {
		dims = &servModel.Dimensions{
			Length: dto.Part.Dimensions.Length,
			Width:  dto.Part.Dimensions.Width,
			Height: dto.Part.Dimensions.Height,
			Weight: dto.Part.Dimensions.Weight,
		}
	}

	var manuf *servModel.Manufacturer
	if dto.Part.Manufacturer != nil {
		manuf = &servModel.Manufacturer{
			Name:    dto.Part.Manufacturer.Name,
			Country: dto.Part.Manufacturer.Country,
			Website: dto.Part.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*servModel.Value, len(dto.Part.Metadata))
	for k, v := range dto.Part.Metadata {
		metadata[k] = valueRepoModelToServModel(v)
	}

	return servModel.PartGetServiceResponse{
		Part: servModel.Part{
			Uuid:          dto.Part.Uuid,
			Name:          dto.Part.Name,
			Description:   dto.Part.Description,
			Price:         dto.Part.Price,
			StockQuantity: dto.Part.StockQuantity,
			Category:      servModel.Category(dto.Part.Category),
			Dimensions:    dims,
			Manufacturer:  manuf,
			Tags:          dto.Part.Tags,
			Metadata:      metadata,
		},
	}
}

func PartListServiceRequestToPartListRepoRequest(dto servModel.PartListServiceRequest) repoModel.PartListRepoRequest {
	categories := []repoModel.Category{}
	if dto.Filter != nil && len(dto.Filter.Categories) > 0 {
		for _, c := range dto.Filter.Categories {
			categories = append(categories, repoModel.Category(c))
		}
	}

	return repoModel.PartListRepoRequest{
		Filter: &repoModel.PartFilterRepo{
			Uuids:                 dto.Filter.Uuids,
			Names:                 dto.Filter.Names,
			Categories:            categories,
			ManufacturerCountries: dto.Filter.ManufacturerCountries,
			Tags:                  dto.Filter.Tags,
		},
	}
}

func PartListRepoResponseToPartListServiceResponse(dto repoModel.PartListRepoResponse) servModel.PartListServiceResponse {
	parts := make([]servModel.Part, 0, len(dto.Parts))
	for _, p := range dto.Parts {
		var dims *servModel.Dimensions
		if p.Dimensions != nil {
			dims = &servModel.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
				Height: p.Dimensions.Height,
				Weight: p.Dimensions.Weight,
			}
		}

		var manuf *servModel.Manufacturer
		if p.Manufacturer != nil {
			manuf = &servModel.Manufacturer{
				Name:    p.Manufacturer.Name,
				Country: p.Manufacturer.Country,
				Website: p.Manufacturer.Website,
			}
		}

		metadata := make(map[string]*servModel.Value, len(p.Metadata))
		for k, v := range p.Metadata {
			metadata[k] = valueRepoModelToServModel(v)
		}

		parts = append(parts, servModel.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      servModel.Category(p.Category),
			Dimensions:    dims,
			Manufacturer:  manuf,
			Tags:          p.Tags,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			Metadata:      metadata,
		})
	}
	return servModel.PartListServiceResponse{Parts: parts}
}

func PartCreateServiceRequestToPartCreateRepoRequest(dto servModel.PartCreateServiceRequest) repoModel.PartCreateRepoRequest {
	var dims *repoModel.Dimensions
	if dto.Dimensions != nil {
		dims = &repoModel.Dimensions{
			Length: dto.Dimensions.Length,
			Width:  dto.Dimensions.Width,
			Height: dto.Dimensions.Height,
			Weight: dto.Dimensions.Weight,
		}
	}

	var manuf *repoModel.Manufacturer
	if dto.Manufacturer != nil {
		manuf = &repoModel.Manufacturer{
			Name:    dto.Manufacturer.Name,
			Country: dto.Manufacturer.Country,
			Website: dto.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*repoModel.Value, len(dto.Metadata))
	for k, v := range dto.Metadata {
		metadata[k] = valueServModelToRepoModel(v)
	}

	return repoModel.PartCreateRepoRequest{
		Part: repoModel.Part{
			Name:          dto.Name,
			Description:   dto.Description,
			Price:         dto.Price,
			StockQuantity: dto.StockQuantity,
			Category:      repoModel.Category(dto.Category),
			Dimensions:    dims,
			Manufacturer:  manuf,
			Tags:          dto.Tags,
			CreatedAt:     dto.CreatedAt,
			UpdatedAt:     dto.UpdatedAt,
			Metadata:      metadata,
		},
	}
}

func valueRepoModelToServModel(v *repoModel.Value) *servModel.Value {
	if v == nil {
		return nil
	}

	if v.StringValue != nil {
		return &servModel.Value{StringValue: v.StringValue}
	}
	if v.Int64Value != nil {
		return &servModel.Value{Int64Value: v.Int64Value}
	}
	if v.DoubleValue != nil {
		return &servModel.Value{DoubleValue: v.DoubleValue}
	}
	if v.BoolValue != nil {
		return &servModel.Value{BoolValue: v.BoolValue}
	}

	return nil
}

func valueServModelToRepoModel(v *servModel.Value) *repoModel.Value {
	if v == nil {
		return nil
	}

	if v.StringValue != nil {
		return &repoModel.Value{StringValue: v.StringValue}
	}
	if v.Int64Value != nil {
		return &repoModel.Value{Int64Value: v.Int64Value}
	}
	if v.DoubleValue != nil {
		return &repoModel.Value{DoubleValue: v.DoubleValue}
	}
	if v.BoolValue != nil {
		return &repoModel.Value{BoolValue: v.BoolValue}
	}

	return nil
}
