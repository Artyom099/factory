package kafka

import (
	"context"

	"github.com/Artyom099/factory/platform/pkg/kafka/consumer"
)

type IConsumer interface {
	Consume(ctx context.Context, handler consumer.MessageHandler) error
}

type IProducer interface {
	Send(ctx context.Context, key, value []byte) error
}
