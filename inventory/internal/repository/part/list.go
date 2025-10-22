package part

import (
	"context"
	"slices"

	"github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/utils"
)

func (r *repository) List(ctx context.Context, dto model.RepoPartFilter) ([]model.RepoPart, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filter := dto
	var parts []model.RepoPart
	for _, part := range r.data {
		if matchPart(part, filter) {
			parts = append(parts, part)
		}
	}

	return parts, nil
}

func matchPart(p model.RepoPart, f model.RepoPartFilter) bool {
	if len(f.Categories) == 0 && len(f.ManufacturerCountries) == 0 && len(f.Names) == 0 && len(f.Tags) == 0 && len(f.Uuids) == 0 {
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
