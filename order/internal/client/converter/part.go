package converter

import (
	"time"

	"github.com/Artyom099/factory/order/internal/service/model"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func PartFilterToProto(filter model.ListPartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func PartFilterToModel(parts []*inventoryV1.Part) model.ListPartsResponseDto {
	outParts := make([]*model.Part, 0, len(parts))
	for _, part := range parts {
		var createdAt, updatedAt time.Time
		if ts := part.GetCreatedAt(); ts != nil {
			createdAt = ts.AsTime()
		}
		if ts := part.GetUpdatedAt(); ts != nil {
			updatedAt = ts.AsTime()
		}

		var dims model.Dimensions
		if d := part.GetDimensions(); d != nil {
			dims = model.Dimensions{
				Length: d.GetLength(),
				Width:  d.GetWidth(),
				Height: d.GetHeight(),
				Weight: d.GetWeight(),
			}
		}

		var manuf model.Manufacturer
		if m := part.GetManufacturer(); m != nil {
			manuf = model.Manufacturer{
				Name:    m.GetName(),
				Country: m.GetCountry(),
				Website: m.GetWebsite(),
			}
		}

		outParts = append(outParts, &model.Part{
			Uuid:          part.GetUuid(),
			Name:          part.GetName(),
			Description:   part.GetDescription(),
			Price:         part.GetPrice(),
			StockQuantity: part.GetStockQuantity(),
			Category:      model.Category(part.GetCategory()),
			Dimensions:    &dims,
			Manufacturer:  &manuf,
			Tags:          part.GetTags(),
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	return model.ListPartsResponseDto{
		Parts: outParts,
	}
}
