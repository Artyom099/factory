package integration

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

func (env *TestEnvironment) GetTestSightingInfo() *inventoryV1.CreatePartRequest {
	return &inventoryV1.CreatePartRequest{
		Name:          "Москва, Красная площадь",
		Description:   "Яркий светящийся объект треугольной формы",
		Price:         50_000.00,
		StockQuantity: 125,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Tags:          []string{"engine"},
	}
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
