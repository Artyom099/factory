package clients

import (
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

const inventoryServerAddress = "localhost:50051"

func CreateInventoryClient() (inventoryV1.InventoryServiceClient, error) {
	// ctx := context.Background()

	connection, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}

	// Создаем gRPC клиент
	return inventoryV1.NewInventoryServiceClient(connection), nil
}
