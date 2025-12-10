package order

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Assemble(ctx context.Context, orderUUID string) error {
	ctx, getOrderSpan := tracing.StartSpan(ctx, "order.ger",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		getOrderSpan.RecordError(err)
		getOrderSpan.End()
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	getOrderSpan.End()

	ctx, updateOrderSpan := tracing.StartSpan(ctx, "order.update_order",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)
	defer updateOrderSpan.End()

	if order.Status == model.OrderStatusPENDINGPAYMENT {
		return model.ErrOrderPendingPayment
	}
	if order.Status == model.OrderStatusASSEMBLED {
		return model.ErrOrderAssembled
	}
	if order.Status == model.OrderStatusCANCELLED {
		return model.ErrOrderCancelled
	}

	order.Status = model.OrderStatusASSEMBLED
	err = s.orderRepository.Update(ctx, order)
	if err != nil {
		updateOrderSpan.RecordError(err)
		return err
	}

	updateOrderSpan.SetAttributes(
		attribute.String("order.status", string(order.Status)),
	)

	return nil
}
