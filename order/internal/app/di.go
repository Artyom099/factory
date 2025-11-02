package app

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderApiV1 "github.com/Artyom099/factory/order/internal/api/order/v1"
	clientGrpc "github.com/Artyom099/factory/order/internal/client/grpc"
	grpcInventoryV1 "github.com/Artyom099/factory/order/internal/client/grpc/inventory/v1"
	grpcPaymentV1 "github.com/Artyom099/factory/order/internal/client/grpc/payment/v1"
	"github.com/Artyom099/factory/order/internal/config"
	"github.com/Artyom099/factory/order/internal/repository"
	orderRepository "github.com/Artyom099/factory/order/internal/repository/order"
	"github.com/Artyom099/factory/order/internal/service"
	orderService "github.com/Artyom099/factory/order/internal/service/order"
	orderV1 "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Artyom099/factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API orderV1.Handler

	orderService service.IOrderService

	orderRepository repository.IOrderRepository

	paymentClient   clientGrpc.IPaymentClient
	inventoryClient clientGrpc.IInventoryClient

	postgresHandle *pgxpool.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OpderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderApiV1.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.IOrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryRepository(ctx),
			d.PaymentClient(ctx),
		)
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.IOrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresHandle(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) PaymentClient(ctx context.Context) clientGrpc.IPaymentClient {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(
			config.AppConfig().Payment.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect payment service: %v", err)
		}

		paymentClient := paymentV1.NewPaymentServiceClient(paymentConn)

		d.paymentClient = grpcPaymentV1.NewClient(paymentClient)
	}

	return d.paymentClient
}

func (d *diContainer) InventoryRepository(ctx context.Context) clientGrpc.IInventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(
			config.AppConfig().Inventory.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			log.Fatalf("failed to connect inventory service: %v", err)
		}

		inventoryClient := inventoryV1.NewInventoryServiceClient(inventoryConn)

		d.inventoryClient = grpcInventoryV1.NewClient(inventoryClient)
	}

	return d.inventoryClient
}

func (d *diContainer) PostgresHandle(ctx context.Context) *pgxpool.Pool {
	if d.postgresHandle == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			log.Fatalf("failed to connect postgres db: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping postgres db: %v\n", err)
		}

		d.postgresHandle = pool
	}

	return d.postgresHandle
}
