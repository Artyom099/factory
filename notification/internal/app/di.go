package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/Artyom099/factory/notification/internal/client/http"
	telegramClient "github.com/Artyom099/factory/notification/internal/client/http/telegram"
	"github.com/Artyom099/factory/notification/internal/config"
	kafkaConverter "github.com/Artyom099/factory/notification/internal/converter/kafka"
	"github.com/Artyom099/factory/notification/internal/converter/kafka/decoder"
	"github.com/Artyom099/factory/notification/internal/service"
	orderAssembledConsumer "github.com/Artyom099/factory/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/Artyom099/factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/Artyom099/factory/notification/internal/service/telegram"
	"github.com/Artyom099/factory/platform/pkg/closer"
	wrappedKafka "github.com/Artyom099/factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Artyom099/factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService               service.INotificationTelegramService
	orderPaidConsumerService      service.INotificationConsumerService
	orderAssembledConsumerService service.INotificationConsumerService

	orderPaidConsumerGroup      sarama.ConsumerGroup
	orderPaidConsumer           wrappedKafka.IConsumer
	orderPaidDecoder            kafkaConverter.IOrderPaidDecoder
	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledConsumer      wrappedKafka.IConsumer
	orderAssembledDecoder       kafkaConverter.IOrderAssembledDecoder

	telegramClient httpClient.ITelegramClient
	telegramBot    *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

// order paid

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.INotificationConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderPaidConsumerService
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderPaid consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderPaid consumer group", func(ctx context.Context) error {
			return d.orderPaidConsumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}

	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.IConsumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
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

// order assembled

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.INotificationConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(
			d.OrderAssembledConsumer(),
			d.OrderAssembledDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create orderAssembled consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka orderAssembled consumer group", func(ctx context.Context) error {
			return d.orderAssembledConsumerGroup.Close()
		})

		d.orderAssembledConsumerGroup = consumerGroup
	}

	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.IConsumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
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

// telegram

func (d *diContainer) TelegramService(ctx context.Context) service.INotificationTelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(
			d.TelegramClient(ctx),
			config.AppConfig().Telegram.ChatID(),
		)
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.ITelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(ctx context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().Telegram.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err.Error()))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}
