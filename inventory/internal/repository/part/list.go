package part

import (
	"context"
	"slices"

	"github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/utils"
)

func (r *repository) List(ctx context.Context, dto model.PartListRepoRequest) (model.PartListRepoResponse, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := dto.Filter
	var parts []model.Part
	for _, part := range r.data {
		if matchPart(part, filter) {
			parts = append(parts, part)
		}
	}

	return model.PartListRepoResponse{Parts: parts}, nil
}

func matchPart(p model.Part, f *model.PartFilterRepo) bool {
	if f == nil {
		return true
	}

	if len(f.Uuids) > 0 && !slices.Contains(f.Uuids, p.Uuid) {
		return false
	}

	if len(f.Names) > 0 && !utils.ContainsInsensitive(f.Names, p.Name) {
		return false
	}

	if len(f.Categories) > 0 && !slices.Contains(f.Categories, p.Category) {
		return false
	}

	if len(f.ManufacturerCountries) > 0 && !utils.ContainsInsensitive(f.ManufacturerCountries, p.Manufacturer.Country) {
		return false
	}

	if len(f.Tags) > 0 && !utils.Intersects(f.Tags, p.Tags) {
		return false
	}

	return true
}
