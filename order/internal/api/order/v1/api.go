package v1

import (
	"github.com/Artyom099/factory/order/internal/service"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

type api struct {
	orderV1.UnimplementedHandler

	orderService service.IOrderService
}

func NewAPI(orderService service.IOrderService) *api {
	return &api{
		orderService: orderService,
	}
}
