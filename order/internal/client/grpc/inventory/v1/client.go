package v1

import (
	def "github.com/Artyom099/factory/order/internal/client/grpc"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

var _ def.IInventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
