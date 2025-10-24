package part

import (
	"context"
	"slices"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
	"github.com/Artyom099/factory/inventory/internal/utils"
)

func (r *repository) List(ctx context.Context, dto model.PartFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	skipFilter := len(dto.Categories) == 0 && len(dto.ManufacturerCountries) == 0 && len(dto.Names) == 0 && len(dto.Tags) == 0 && len(dto.Uuids) == 0

	var parts []model.Part
	for _, part := range r.data {
		if skipFilter {
			parts = append(parts, converter.ToModelPart(part))
		} else if matchPart(part, converter.ToRepoPartFilter(dto)) {
			parts = append(parts, converter.ToModelPart(part))
		}
	}

	return parts, nil
}

func matchPart(p repoModel.RepoPart, f repoModel.RepoPartFilter) bool {
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
