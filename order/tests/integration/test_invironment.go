package integration

import (
	"context"

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

func (env *TestEnvironment) ClearDB(ctx context.Context) error {
	// databaseName := os.Getenv("POSTGRES_DB")
	// if databaseName == "" {
	// 	databaseName = "order-service" // fallback значение
	// }

	env.Postgres.Client()

	// Database(databaseName).
	// Collection(partsCollectionName).
	// DeleteMany(ctx, bson.M{})

	// if err != nil {
	// 	return err
	// }

	return nil
}
