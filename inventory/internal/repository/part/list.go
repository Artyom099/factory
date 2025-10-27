package part

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (r *repository) List(ctx context.Context, dto model.PartFilter) ([]model.Part, error) {
	filter := buildMongoFilter(converter.ToRepoPartFilter(dto))

	cursor, err := r.collection.Find(ctx, filter, options.Find())
	if err != nil {
		return nil, fmt.Errorf("failed to query parts: %w", err)
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", cerr)
		}
	}()

	var repoParts []repoModel.RepoPart
	if err := cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	var parts []model.Part
	for _, part := range repoParts {
		parts = append(parts, converter.ToModelPart(part))
	}

	return parts, nil
}

func buildMongoFilter(f repoModel.RepoPartFilter) bson.M {
	filter := bson.M{}

	if len(f.Uuids) > 0 {
		filter["uuid"] = bson.M{"$in": f.Uuids}
	}
	if len(f.Names) > 0 {
		filter["name"] = bson.M{"$in": f.Names}
	}
	if len(f.Categories) > 0 {
		filter["category"] = bson.M{"$in": f.Categories}
	}
	if len(f.ManufacturerCountries) > 0 {
		filter["manufacturer.country"] = bson.M{"$in": f.ManufacturerCountries}
	}
	if len(f.Tags) > 0 {
		filter["tags"] = bson.M{"$in": f.Tags}
	}

	return filter
}
