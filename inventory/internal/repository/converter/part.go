package converter

import (
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

func ToModelPart(dto repoModel.RepoPart) servModel.Part {
	var dimensionModel servModel.Dimensions
	if dto.Dimensions != nil {
		dimensionModel = servModel.Dimensions{
			Length: dto.Dimensions.Length,
			Width:  dto.Dimensions.Width,
			Height: dto.Dimensions.Height,
			Weight: dto.Dimensions.Weight,
		}
	}

	var manufacturerModel servModel.Manufacturer
	if dto.Manufacturer != nil {
		manufacturerModel = servModel.Manufacturer{
			Name:    dto.Manufacturer.Name,
			Country: dto.Manufacturer.Country,
			Website: dto.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*servModel.Value, len(dto.Metadata))
	for k, v := range dto.Metadata {
		metadata[k] = toModelPartValue(v)
	}

	return servModel.Part{
		Uuid:          dto.Uuid,
		Name:          dto.Name,
		Description:   dto.Description,
		Price:         dto.Price,
		StockQuantity: dto.StockQuantity,
		Category:      servModel.Category(dto.Category),
		Dimensions:    &dimensionModel,
		Manufacturer:  &manufacturerModel,
		Tags:          dto.Tags,
		Metadata:      metadata,
	}
}

func ToRepoPartFilter(dto servModel.PartFilter) repoModel.RepoPartFilter {
	categories := []repoModel.Category{}
	if len(dto.Categories) > 0 {
		for _, c := range dto.Categories {
			categories = append(categories, repoModel.Category(c))
		}
	}

	return repoModel.RepoPartFilter{
		Uuids:                 dto.Uuids,
		Names:                 dto.Names,
		Categories:            categories,
		ManufacturerCountries: dto.ManufacturerCountries,
		Tags:                  dto.Tags,
	}
}

func ToModelListParts(dto []repoModel.RepoPart) []servModel.Part {
	parts := make([]servModel.Part, 0, len(dto))
	for _, p := range dto {
		var dimensionModel servModel.Dimensions
		if p.Dimensions != nil {
			dimensionModel = servModel.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
				Height: p.Dimensions.Height,
				Weight: p.Dimensions.Weight,
			}
		}

		var manufacturerModel servModel.Manufacturer
		if p.Manufacturer != nil {
			manufacturerModel = servModel.Manufacturer{
				Name:    p.Manufacturer.Name,
				Country: p.Manufacturer.Country,
				Website: p.Manufacturer.Website,
			}
		}

		metadata := make(map[string]*servModel.Value, len(p.Metadata))
		for k, v := range p.Metadata {
			metadata[k] = toModelPartValue(v)
		}

		parts = append(parts, servModel.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      servModel.Category(p.Category),
			Dimensions:    &dimensionModel,
			Manufacturer:  &manufacturerModel,
			Tags:          p.Tags,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
			Metadata:      metadata,
		})
	}

	return parts
}

func ToRepoPart(dto servModel.Part) repoModel.RepoPart {
	var dimensionModel repoModel.Dimensions
	if dto.Dimensions != nil {
		dimensionModel = repoModel.Dimensions{
			Length: dto.Dimensions.Length,
			Width:  dto.Dimensions.Width,
			Height: dto.Dimensions.Height,
			Weight: dto.Dimensions.Weight,
		}
	}

	var manufacturerModel repoModel.Manufacturer
	if dto.Manufacturer != nil {
		manufacturerModel = repoModel.Manufacturer{
			Name:    dto.Manufacturer.Name,
			Country: dto.Manufacturer.Country,
			Website: dto.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*repoModel.Value, len(dto.Metadata))
	for k, v := range dto.Metadata {
		metadata[k] = toRepoPartValue(v)
	}

	return repoModel.RepoPart{
		Uuid:          dto.Uuid,
		Name:          dto.Name,
		Description:   dto.Description,
		Price:         dto.Price,
		StockQuantity: dto.StockQuantity,
		Category:      repoModel.Category(dto.Category),
		Dimensions:    &dimensionModel,
		Manufacturer:  &manufacturerModel,
		Tags:          dto.Tags,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
		Metadata:      metadata,
	}
}

func toModelPartValue(v *repoModel.Value) *servModel.Value {
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

func toRepoPartValue(v *servModel.Value) *repoModel.Value {
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
