package v1

import (
	"github.com/Artyom099/factory/payment/internal/service"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

type api struct {
	paymentV1.UnimplementedPaymentServiceServer

	paymentService service.IPaymentService
}

func NewAPI(paymentService service.IPaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
