package app

import (
	"context"

	paymentApiV1 "github.com/Artyom099/factory/payment/internal/api/payment/v1"
	"github.com/Artyom099/factory/payment/internal/repository"
	paymentRepository "github.com/Artyom099/factory/payment/internal/repository/payment"
	"github.com/Artyom099/factory/payment/internal/service"
	paymentService "github.com/Artyom099/factory/payment/internal/service/payment"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API paymentV1.PaymentServiceServer

	paymentService service.IPaymentService

	paymentRepository repository.IPaymentRepository
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentApiV1.NewAPI(d.PartService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PartService(ctx context.Context) service.IPaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService(d.PartRepository(ctx))
	}

	return d.paymentService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.IPaymentRepository {
	if d.paymentRepository == nil {
		d.paymentRepository = paymentRepository.NewRepository()
	}

	return d.paymentRepository
}
