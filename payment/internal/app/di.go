package app

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	paymentApiV1 "github.com/Artyom099/factory/payment/internal/api/payment/v1"
	"github.com/Artyom099/factory/payment/internal/config"
	"github.com/Artyom099/factory/payment/internal/repository"
	paymentRepository "github.com/Artyom099/factory/payment/internal/repository/payment"
	"github.com/Artyom099/factory/payment/internal/service"
	paymentService "github.com/Artyom099/factory/payment/internal/service/payment"
	grpcAuth "github.com/Artyom099/factory/platform/pkg/middleware/grpc"
	"github.com/Artyom099/factory/platform/pkg/tracing"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API      paymentV1.PaymentServiceServer
	paymentService    service.IPaymentService
	paymentRepository repository.IPaymentRepository

	iamClient grpcAuth.IAMClient
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentApiV1.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(ctx context.Context) service.IPaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService(d.PaymentRepository(ctx))
	}

	return d.paymentService
}

func (d *diContainer) PaymentRepository(ctx context.Context) repository.IPaymentRepository {
	if d.paymentRepository == nil {
		d.paymentRepository = paymentRepository.NewRepository()
	}

	return d.paymentRepository
}

func (d *diContainer) IAMClient(ctx context.Context) grpcAuth.IAMClient {
	if d.iamClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IamCLient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(tracing.UnaryClientInterceptor("iam-service")),
		)
		if err != nil {
			log.Fatalf("failed to connect IAM service: %v", err)
		}
		d.iamClient = authV1.NewAuthServiceClient(conn)
	}

	return d.iamClient
}
