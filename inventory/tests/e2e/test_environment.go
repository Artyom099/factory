package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (env *TestEnvironment) GetCreatePartRequest() *inventoryV1.CreatePartRequest {
	return &inventoryV1.CreatePartRequest{
		Name:          "Right Engine",
		Description:   "Turbina Engine",
		Price:         50_000.00,
		StockQuantity: 125,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Tags:          []string{"engine, right"},
	}
}

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	uuid := gofakeit.UUID()
	now := time.Now()

	part := bson.M{
		"uuid":           uuid,
		"name":           gofakeit.Name(),
		"description":    gofakeit.ProductDescription(),
		"price":          gofakeit.Price(10, 20_000),
		"category":       inventoryV1.Category_CATEGORY_ENGINE,
		"stock_quantity": gofakeit.Uint8(),
		"dimensions": bson.M{
			"width":  gofakeit.Float64(),
			"height": gofakeit.Float64(),
			"length": gofakeit.Float64(),
			"weight": gofakeit.Float64(),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.AppName(),
			"country": gofakeit.Country(),
		},
		"tags": []string{gofakeit.Adjective(), gofakeit.Adjective()},
		// "metadata": bson.M{
		// 	"key1": gofakeit.BeerName(),
		// },
		"created_at": primitive.NewDateTimeFromTime(now),
	}

	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, part)
	if err != nil {
		return "", err
	}

	return uuid, nil
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
