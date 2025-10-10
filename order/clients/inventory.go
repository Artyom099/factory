package clients

import (
	"log"

	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const inventoryServerAddress = "localhost:50051"

func CreateInventoryClient() (inventoryV1.InventoryV1ServiceClient, error) {
	// ctx := context.Background()

	connection, err := grpc.NewClient(
		inventoryServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return nil, err
	}
	defer func() {
		if cerr := connection.Close(); cerr != nil {
			log.Printf("failed to close connect: %v", cerr)
		}
	}()

	// Создаем gRPC клиент
	return inventoryV1.NewInventoryV1ServiceClient(connection), nil
}
