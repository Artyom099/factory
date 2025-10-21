package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

// converter api - service

func PartGetApiRequestToPartGetServiceRequest(dto *inventoryV1.GetPartRequest) servModel.PartGetServiceRequest {
	return servModel.PartGetServiceRequest{Uuid: dto.GetUuid()}
}

func PartGetServiceResponseToPartGetApiResponse(dto servModel.PartGetServiceResponse) *inventoryV1.GetPartResponse {
	if dto.Part.Uuid == "" {
		return nil
	}

	var dims *inventoryV1.Dimensions
	if dto.Part.Dimensions != nil {
		dims = &inventoryV1.Dimensions{
			Length: dto.Part.Dimensions.Length,
			Width:  dto.Part.Dimensions.Width,
			Height: dto.Part.Dimensions.Height,
			Weight: dto.Part.Dimensions.Weight,
		}
	}

	var manuf *inventoryV1.Manufacturer
	if dto.Part.Manufacturer != nil {
		manuf = &inventoryV1.Manufacturer{
			Name:    dto.Part.Manufacturer.Name,
			Country: dto.Part.Manufacturer.Country,
			Website: dto.Part.Manufacturer.Website,
		}
	}

	metadata := make(map[string]*inventoryV1.Value, len(dto.Part.Metadata))
	for k, v := range dto.Part.Metadata {
		metadata[k] = valueServModelToProto(v)
	}

	part := &inventoryV1.Part{
		Uuid:          dto.Part.Uuid,
		Name:          dto.Part.Name,
		Description:   dto.Part.Description,
		Price:         dto.Part.Price,
		StockQuantity: dto.Part.StockQuantity,
		Category:      inventoryV1.Category(dto.Part.Category),
		Dimensions:    dims,
		Manufacturer:  manuf,
		Tags:          dto.Part.Tags,
		Metadata:      metadata,
		CreatedAt:     timestamppb.New(dto.Part.CreatedAt),
		UpdatedAt:     timestamppb.New(dto.Part.UpdatedAt),
	}

	return &inventoryV1.GetPartResponse{Part: part}
}

func PartListApiRequestToPartGetServiceRequest(dto *inventoryV1.ListPartsRequest) servModel.PartListServiceRequest {
	if dto == nil || dto.GetFilter() == nil {
		return servModel.PartListServiceRequest{}
	}

	f := dto.GetFilter()
	svcFilter := &servModel.PartsFilterService{
		Uuids:                 f.GetUuids(),
		Names:                 f.GetNames(),
		Categories:            nil,
		ManufacturerCountries: f.GetManufacturerCountries(),
		Tags:                  f.GetTags(),
	}

	if len(f.GetCategories()) > 0 {
		svcFilter.Categories = make([]servModel.Category, 0, len(f.GetCategories()))
		for _, c := range f.GetCategories() {
			svcFilter.Categories = append(svcFilter.Categories, servModel.Category(c))
		}
	}

	return servModel.PartListServiceRequest{Filter: svcFilter}
}

func PartListServiceResponseToPartGetApiResponse(dto servModel.PartListServiceResponse) *inventoryV1.ListPartsResponse {
	parts := make([]*inventoryV1.Part, 0, len(dto.Parts))
	for _, p := range dto.Parts {
		var dims *inventoryV1.Dimensions
		if p.Dimensions != nil {
			dims = &inventoryV1.Dimensions{
				Length: p.Dimensions.Length,
				Width:  p.Dimensions.Width,
				Height: p.Dimensions.Height,
				Weight: p.Dimensions.Weight,
			}
		}

		var manuf *inventoryV1.Manufacturer
		if p.Manufacturer != nil {
			manuf = &inventoryV1.Manufacturer{
				Name:    p.Manufacturer.Name,
				Country: p.Manufacturer.Country,
				Website: p.Manufacturer.Website,
			}
		}

		metadata := make(map[string]*inventoryV1.Value, len(p.Metadata))
		for k, v := range p.Metadata {
			metadata[k] = valueServModelToProto(v)
		}

		parts = append(parts, &inventoryV1.Part{
			Uuid:          p.Uuid,
			Name:          p.Name,
			Description:   p.Description,
			Price:         p.Price,
			StockQuantity: p.StockQuantity,
			Category:      inventoryV1.Category(p.Category),
			Dimensions:    dims,
			Manufacturer:  manuf,
			Tags:          p.Tags,
			Metadata:      metadata,
			CreatedAt:     timestamppb.New(p.CreatedAt),
			UpdatedAt:     timestamppb.New(p.UpdatedAt),
		})
	}

	return &inventoryV1.ListPartsResponse{Parts: parts}
}

func PartCreateApiRequestToPartCreateServiceRequest(dto *inventoryV1.CreatePartRequest) servModel.PartCreateServiceRequest {
	var dims *servModel.Dimensions
	if d := dto.GetDimensions(); d != nil {
		dims = &servModel.Dimensions{
			Length: d.Length,
			Width:  d.Width,
			Height: d.Height,
			Weight: d.Weight,
		}
	}

	var manuf *servModel.Manufacturer
	if m := dto.GetManufacturer(); m != nil {
		manuf = &servModel.Manufacturer{
			Name:    m.Name,
			Country: m.Country,
			Website: m.Website,
		}
	}

	metadata := make(map[string]*servModel.Value, len(dto.GetMetadata()))
	for k, v := range dto.GetMetadata() {
		metadata[k] = valueProtoToServModel(v)
	}

	return servModel.PartCreateServiceRequest{
		Name:          dto.GetName(),
		Description:   dto.GetDescription(),
		Price:         dto.GetPrice(),
		StockQuantity: dto.GetStockQuantity(),
		Category:      servModel.Category(dto.GetCategory()),
		Dimensions:    dims,
		Manufacturer:  manuf,
		Tags:          dto.GetTags(),
		Metadata:      metadata,
	}
}

func PartCreateServiceResponseToPartCreateApiResponse(dto servModel.PartCreateServiceResponse) *inventoryV1.CreatePartResponse {
	return &inventoryV1.CreatePartResponse{Uuid: dto.Uuid}
}

func valueProtoToServModel(v *inventoryV1.Value) *servModel.Value {
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

func valueServModelToProto(v *servModel.Value) *inventoryV1.Value {
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
