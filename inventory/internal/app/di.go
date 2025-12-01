package app

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	partApiV1 "github.com/Artyom099/factory/inventory/internal/api/part/v1"
	"github.com/Artyom099/factory/inventory/internal/config"
	"github.com/Artyom099/factory/inventory/internal/repository"
	partRepository "github.com/Artyom099/factory/inventory/internal/repository/part"
	"github.com/Artyom099/factory/inventory/internal/service"
	partService "github.com/Artyom099/factory/inventory/internal/service/part"
	"github.com/Artyom099/factory/platform/pkg/closer"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	partService service.IPartService

	partRepository repository.IPartRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database

	iamClient authV1.AuthServiceClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = partApiV1.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.IPartService {
	if d.partService == nil {
		d.partService = partService.NewService(d.PartRepository(ctx))
	}

	return d.partService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.IPartRepository {
	if d.partRepository == nil {
		d.partRepository = partRepository.NewRepository(ctx, d.MongoDBHandle(ctx))
	}

	return d.partRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) IAMClient(ctx context.Context) authV1.AuthServiceClient {
	if d.iamClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IamCLient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect IAM service: %v", err)
		}
		d.iamClient = authV1.NewAuthServiceClient(conn)
	}

	return d.iamClient
}
