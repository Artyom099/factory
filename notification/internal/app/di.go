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
	orderPaidConsumer "github.com/Artyom099/factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/Artyom099/factory/notification/internal/service/telegram"
	"github.com/Artyom099/factory/platform/pkg/closer"
	wrappedKafka "github.com/Artyom099/factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Artyom099/factory/platform/pkg/kafka/consumer"
	"github.com/Artyom099/factory/platform/pkg/logger"
	kafkaMiddleware "github.com/Artyom099/factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService             service.INotificationTelegramService
	notificationConsumerService service.INotificationConsumerService

	consumerGroup          sarama.ConsumerGroup
	orderPaidConsumer      wrappedKafka.IConsumer
	orderPaidDecoder       kafkaConverter.IOrderPaidDecoder
	orderAssembledConsumer wrappedKafka.IConsumer
	orderAssembledDecoder  kafkaConverter.IOrderAssembledDecoder

	telegramClient httpClient.ITelegramClient
	telegramBot    *bot.Bot
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) NotificationConsumerService(ctx context.Context) service.INotificationConsumerService {
	if d.notificationConsumerService == nil {
		d.notificationConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(ctx),
		)
	}

	return d.notificationConsumerService
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
