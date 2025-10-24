package v1

import (
	"github.com/Artyom099/factory/inventory/internal/service"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.IPartService
}

func NewAPI(partService service.IPartService) *api {
	return &api{
		partService: partService,
	}
}
