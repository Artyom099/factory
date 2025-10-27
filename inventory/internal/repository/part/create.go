package part

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (r *repository) Create(ctx context.Context, dto model.Part) (string, error) {
	part := converter.ToRepoPart(dto)
	uuid := uuid.New().String()
	part.Uuid = uuid

	if part.CreatedAt.IsZero() {
		*part.CreatedAt = time.Now()
	}

	res, err := r.collection.InsertOne(ctx, part)
	if err != nil {
		return "", model.ErrInternalError
	}

	var createdPart repoModel.RepoPart
	err = r.collection.FindOne(ctx, bson.M{"_id": res.InsertedID}).Decode(&createdPart)
	if err != nil {
		log.Printf("repo_err: %v", err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", repoModel.ErrPartNotFound
		}

		return "", repoModel.ErrInternalError
	}

	return createdPart.Uuid, nil
}
