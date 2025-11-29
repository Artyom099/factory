package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/Artyom099/factory/notification/internal/model"
	eventsV1 "github.com/Artyom099/factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderPaidInEvent, error) {
	var pb eventsV1.OrderPaidEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidInEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidInEvent{
		EventUUID:       pb.EventUuid,
		OrderUUID:       pb.OrderUuid,
		UserUUID:        pb.UserUuid,
		PaymentMethod:   pb.PaymentMethod.String(),
		TransactionUUID: pb.TransactionUuid,
	}, nil
}
