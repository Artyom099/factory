package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (model.Part, error) {
	var part repoModel.RepoPart
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, repoModel.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return converter.ToModelPart(part), nil
}
