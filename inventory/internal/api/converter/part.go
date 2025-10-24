package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func ToApiPart(dto servModel.Part) *inventoryV1.GetPartResponse {
	if dto.Uuid == "" {
		return nil
	}

	var dimensionApi inventoryV1.Dimensions
	if dto.Dimensions != nil {
		dimensionApi = inventoryV1.Dimensions{
			Length: dto.Dimensions.Length,
			Width:  dto.Dimensions.Width,
			Height: dto.Dimensions.Height,
			Weight: dto.Dimensions.Weight,
		}
	}

	var manufacturerApi inventoryV1.Manufacturer
	if dto.Manufacturer != nil {
		manufacturerApi = inventoryV1.Manufacturer{
			Name:    dto.Manufacturer.Name,
			Country: dto.Manufacturer.Country,
			Website: dto.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*inventoryV1.Value, len(dto.Metadata))
	for k, v := range dto.Metadata {
		metadata[k] = toApiPartValue(v)
	}

	part := &inventoryV1.Part{
		Uuid:          dto.Uuid,
		Name:          dto.Name,
		Description:   dto.Description,
		Price:         dto.Price,
		StockQuantity: dto.StockQuantity,
		Category:      inventoryV1.Category(dto.Category),
		Dimensions:    &dimensionApi,
		Manufacturer:  &manufacturerApi,
		Tags:          dto.Tags,
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(dto.CreatedAt),
		UpdatedAt:     timestamppb.New(dto.UpdatedAt),
	}

	return &inventoryV1.GetPartResponse{Part: part}
}

func ToModelPartFilter(dto *inventoryV1.ListPartsRequest) servModel.PartFilter {
	if dto == nil || dto.GetFilter() == nil {
		return servModel.PartFilter{}
	}

	f := dto.GetFilter()
	filter := servModel.PartFilter{
		Uuids:                 f.GetUuids(),
		Names:                 f.GetNames(),
		Categories:            nil,
		ManufacturerCountries: f.GetManufacturerCountries(),
		Tags:                  f.GetTags(),
	}

	if len(f.GetCategories()) > 0 {
		filter.Categories = make([]servModel.Category, 0, len(f.GetCategories()))
		for _, c := range f.GetCategories() {
			filter.Categories = append(filter.Categories, servModel.Category(c))
		}
	}

	return filter
}

func ToApiListParts(dto []servModel.Part) *inventoryV1.ListPartsResponse {
	parts := make([]*inventoryV1.Part, 0, len(dto))
	for _, p := range dto {
		var dimensionApi inventoryV1.Dimensions
		if p.Dimensions != nil {
			dimensionApi = inventoryV1.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
				Height: p.Dimensions.Height,
				Weight: p.Dimensions.Weight,
			}
		}

		var manufacturerApi inventoryV1.Manufacturer
		if p.Manufacturer != nil {
			manufacturerApi = inventoryV1.Manufacturer{
				Name:    p.Manufacturer.Name,
				Country: p.Manufacturer.Country,
				Website: p.Manufacturer.Website,
			}
		}

		metadata := make(map[string]*inventoryV1.Value, len(p.Metadata))
		for k, v := range p.Metadata {
			metadata[k] = toApiPartValue(v)
		}

		parts = append(parts, &inventoryV1.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      inventoryV1.Category(p.Category),
			Dimensions:    &dimensionApi,
			Manufacturer:  &manufacturerApi,
			Tags:          p.Tags,
			Metadata:      metadata,
			CreatedAt:     timestamppb.New(p.CreatedAt),
			UpdatedAt:     timestamppb.New(p.UpdatedAt),
		})
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}
}

func ToModelPart(dto *inventoryV1.CreatePartRequest) servModel.Part {
	var dimensionsModel servModel.Dimensions
	if d := dto.GetDimensions(); d != nil {
		dimensionsModel = servModel.Dimensions{
			Length: d.Length,
			Width:  d.Width,
			Height: d.Height,
			Weight: d.Weight,
		}
	}

	var manufacturerModel servModel.Manufacturer
	if m := dto.GetManufacturer(); m != nil {
		manufacturerModel = servModel.Manufacturer{
			Name:    m.Name,
			Country: m.Country,
			Website: m.Website,
		}
	}

	metadata := make(map[string]*servModel.Value, len(dto.GetMetadata()))
	for k, v := range dto.GetMetadata() {
		metadata[k] = toModelPartValue(v)
	}

	return servModel.Part{
		Name:          dto.GetName(),
		Description:   dto.GetDescription(),
		Price:         dto.GetPrice(),
		StockQuantity: dto.GetStockQuantity(),
		Category:      servModel.Category(dto.GetCategory()),
		Dimensions:    &dimensionsModel,
		Manufacturer:  &manufacturerModel,
		Tags:          dto.GetTags(),
		Metadata:      metadata,
	}
}

func toModelPartValue(v *inventoryV1.Value) *servModel.Value {
	if v == nil {
		return nil
	}

	switch kind := v.Kind.(type) {
	case *inventoryV1.Value_StringValue:
		return &servModel.Value{StringValue: &kind.StringValue}
	case *inventoryV1.Value_Int64Value:
		return &servModel.Value{Int64Value: &kind.Int64Value}
	case *inventoryV1.Value_DoubleValue:
		return &servModel.Value{DoubleValue: &kind.DoubleValue}
	case *inventoryV1.Value_BoolValue:
		return &servModel.Value{BoolValue: &kind.BoolValue}
	default:
		return nil
	}
}

func toApiPartValue(v *servModel.Value) *inventoryV1.Value {
	if v == nil {
		return nil
	}

	if v.StringValue != nil {
		return &inventoryV1.Value{Kind: &inventoryV1.Value_StringValue{StringValue: *v.StringValue}}
	}
	if v.Int64Value != nil {
		return &inventoryV1.Value{Kind: &inventoryV1.Value_Int64Value{Int64Value: *v.Int64Value}}
	}
	if v.DoubleValue != nil {
		return &inventoryV1.Value{Kind: &inventoryV1.Value_DoubleValue{DoubleValue: *v.DoubleValue}}
	}
	if v.BoolValue != nil {
		return &inventoryV1.Value{Kind: &inventoryV1.Value_BoolValue{BoolValue: *v.BoolValue}}
	}

	return nil
}
