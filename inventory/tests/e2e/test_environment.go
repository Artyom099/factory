package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

type InsertedPart struct {
	Uuid          string                       `bson:"uuid"`
	Name          string                       `bson:"name"`
	Description   string                       `bson:"description"`
	Price         float64                      `bson:"price"`
	Category      inventoryV1.Category         `bson:"category"`
	StockQuantity int64                        `bson:"stock_quantity"`
	Dimensions    *inventoryV1.Dimensions      `bson:"dimensions"`
	Manufacturer  *inventoryV1.Manufacturer    `bson:"manufacturer"`
	Tags          []string                     `bson:"tags"`
	Metadata      map[string]map[string]string `bson:"metadata"`
}

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

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (*InsertedPart, error) {
	uuid := gofakeit.UUID()

	part := bson.M{
		"uuid":           uuid,
		"name":           gofakeit.Product().Name,
		"description":    gofakeit.Product().Description,
		"price":          gofakeit.Price(10, 20_000),
		"category":       inventoryV1.Category(gofakeit.Number(0, 4)), //nolint:gosec
		"stock_quantity": gofakeit.Int64(),
		"dimensions": bson.M{
			"width":  gofakeit.Float64Range(0, 100_000),
			"height": gofakeit.Float64Range(0, 100_000),
			"length": gofakeit.Float64Range(0, 100_000),
			"weight": gofakeit.Float64Range(0, 100_000),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
		},
		"tags": []string{gofakeit.Adjective(), gofakeit.Adjective()},
		"metadata": bson.M{
			"material": bson.M{
				"string_value": gofakeit.BeerName(),
			},
		},
		"created_at": primitive.NewDateTimeFromTime(time.Now()),
	}

	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	result, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, part)
	if err != nil {
		return &InsertedPart{}, err
	}

	var inserted InsertedPart
	err = env.Mongo.Client().
		Database(databaseName).
		Collection(partsCollectionName).
		FindOne(ctx, bson.M{"_id": result.InsertedID}).
		Decode(&inserted)
	if err != nil {
		return &InsertedPart{}, err
	}

	return &inserted, nil
}

func AssertPartsEqual(expected *InsertedPart, got *inventoryV1.Part) {
	Expect := gomega.Expect
	Equal := gomega.Equal
	BeTrue := gomega.BeTrue

	Expect(got.Uuid).To(Equal(expected.Uuid))
	Expect(got.Name).To(Equal(expected.Name))
	Expect(got.Description).To(Equal(expected.Description))
	Expect(got.Price).To(Equal(expected.Price))
	Expect(got.Category).To(Equal(expected.Category))
	Expect(got.StockQuantity).To(Equal(expected.StockQuantity))
	Expect(got.Dimensions.Width).To(Equal(expected.Dimensions.Width))
	Expect(got.Dimensions.Height).To(Equal(expected.Dimensions.Height))
	Expect(got.Dimensions.Length).To(Equal(expected.Dimensions.Length))
	Expect(got.Dimensions.Weight).To(Equal(expected.Dimensions.Weight))
	Expect(got.Manufacturer.Name).To(Equal(expected.Manufacturer.Name))
	Expect(got.Manufacturer.Country).To(Equal(expected.Manufacturer.Country))
	Expect(got.Tags).To(Equal(expected.Tags))

	if expected.Metadata != nil && got.Metadata != nil {
		for key, expVal := range expected.Metadata {
			gotVal, ok := got.Metadata[key]
			Expect(ok).To(BeTrue(), "missing metadata key %s", key)

			switch val := any(expVal).(type) {
			case string:
				Expect(gotVal.GetStringValue()).
					To(Equal(val), "mismatch in metadata[%s]", key)
			case map[string]string:
				Expect(gotVal.GetStringValue()).
					To(Equal(val["string_value"]), "mismatch in metadata[%s]", key)
			default:
				Expect(false).To(BeTrue(), "unsupported metadata type for %s", key)
			}
		}
	}
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
