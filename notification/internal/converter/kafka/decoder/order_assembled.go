package decoder

import (
	"fmt"

	"github.com/gogo/protobuf/proto"

	"github.com/Artyom099/factory/notification/internal/model"
	eventsV1 "github.com/Artyom099/factory/shared/pkg/proto/events/v1"
)

type orderAssembledDecoder struct{}

func NewOrderAssembledDecoder() *orderAssembledDecoder {
	return &orderAssembledDecoder{}
}

func (d *orderAssembledDecoder) Decode(data []byte) (model.OrderAssembledInEvent, error) {
	var pb eventsV1.ShipAssembledEvent
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
