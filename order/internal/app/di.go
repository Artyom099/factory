package app

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	kafkaConverter "github.com/Artyom099/factory/order/internal/api/converter/kafka"
	"github.com/Artyom099/factory/order/internal/api/converter/kafka/decoder"
	orderApiV1 "github.com/Artyom099/factory/order/internal/api/order/v1"
	clientGrpc "github.com/Artyom099/factory/order/internal/client/grpc"
	grpcInventoryV1 "github.com/Artyom099/factory/order/internal/client/grpc/inventory/v1"
	grpcPaymentV1 "github.com/Artyom099/factory/order/internal/client/grpc/payment/v1"
	"github.com/Artyom099/factory/order/internal/config"
	"github.com/Artyom099/factory/order/internal/repository"
	orderRepository "github.com/Artyom099/factory/order/internal/repository/order"
	"github.com/Artyom099/factory/order/internal/service"
	orderConsumer "github.com/Artyom099/factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/Artyom099/factory/order/internal/service/order"
	orderProducer "github.com/Artyom099/factory/order/internal/service/producer/order_producer"
	"github.com/Artyom099/factory/platform/pkg/closer"
	wrappedKafka "github.com/Artyom099/factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Artyom099/factory/platform/pkg/kafka/producer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Artyom099/factory/platform/pkg/middleware/kafka"
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

	orderProducerService service.IOrderProducerService
	orderConsumerService service.IOrderConsumerService

	consumerGroup          sarama.ConsumerGroup
	orderAssembledConsumer wrappedKafka.IConsumer
	orderAssembledDecoder  kafkaConverter.IOrderAssembledDecoder

	syncProducer      sarama.SyncProducer
	orderPaidProducer wrappedKafka.IProducer
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
			d.OrderProducerService(),
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

// Kafka

func (d *diContainer) OrderProducerService() service.IOrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderPaidProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService(ctx context.Context) service.IOrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(
			d.OrderService(ctx),
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
		)
	}

	return d.orderConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.IConsumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.IOrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer() wrappedKafka.IProducer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderPaidProducer
}
