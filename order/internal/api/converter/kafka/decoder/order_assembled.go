package decoder

import (
	"fmt"

	"google.golang.org/protobuf/proto"

	"github.com/Artyom099/factory/order/internal/service/model"
	eventsV1 "github.com/Artyom099/factory/shared/pkg/proto/events/v1"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.OrderAssembledInEvent, error) {
	var pb eventsV1.OrderAssembledEvent
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderAssembledInEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderAssembledInEvent{
		EventUUID:    pb.EventUuid,
		OrderUUID:    pb.OrderUuid,
		UserUUID:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
