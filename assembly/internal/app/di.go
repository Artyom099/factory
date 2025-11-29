package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/Artyom099/factory/assembly/internal/config"
	kafkaConverter "github.com/Artyom099/factory/assembly/internal/converter/kafka"
	"github.com/Artyom099/factory/assembly/internal/converter/kafka/decoder"
	"github.com/Artyom099/factory/assembly/internal/service"
	assemblyService "github.com/Artyom099/factory/assembly/internal/service/assembly"
	orderConsumer "github.com/Artyom099/factory/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/Artyom099/factory/assembly/internal/service/producer/order_producer"
	"github.com/Artyom099/factory/platform/pkg/closer"
	wrappedKafka "github.com/Artyom099/factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Artyom099/factory/platform/pkg/kafka/producer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Artyom099/factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	assemblyService         service.IAssemblyService
	assemblyProducerService service.IAssemblyProducerService
	assemblyConsumerService service.IAssemblyConsumerService

	consumerGroup     sarama.ConsumerGroup
	orderPaidConsumer wrappedKafka.IConsumer
	orderPaidDecoder  kafkaConverter.IOrderPaidDecoder

	syncProducer           sarama.SyncProducer
	orderAssembledProducer wrappedKafka.IProducer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) AssemblyProducerService() service.IAssemblyProducerService {
	if d.assemblyProducerService == nil {
		d.assemblyProducerService = orderProducer.NewService(d.OrderAssembledProducer())
	}

	return d.assemblyProducerService
}

func (d *diContainer) AssemblyConsumerService() service.IAssemblyConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = orderConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.AssemblyService(),
		)
	}

	return d.assemblyConsumerService
}

func (d *diContainer) AssemblyService() service.IAssemblyService {
	if d.assemblyService == nil {
		d.assemblyService = assemblyService.NewService(d.AssemblyProducerService())
	}

	return d.assemblyService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
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

func (d *diContainer) OrderPaidConsumer() wrappedKafka.IConsumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.IOrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
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

func (d *diContainer) OrderAssembledProducer() wrappedKafka.IProducer {
	if d.orderAssembledProducer == nil {
		d.orderAssembledProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssembledProducer
}
