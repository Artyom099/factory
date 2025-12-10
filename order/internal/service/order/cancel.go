package order

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	ctx, getOrderSpan := tracing.StartSpan(ctx, "order.get",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		getOrderSpan.RecordError(err)
		getOrderSpan.End()
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	getOrderSpan.End()

	ctx, cancelOrderSpan := tracing.StartSpan(ctx, "order.cancel",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)
	defer cancelOrderSpan.End()

	if order.Status == model.OrderStatusPAID || order.Status == model.OrderStatusCANCELLED {
		return model.ErrConflict
	}

	if order.Status == model.OrderStatusPENDINGPAYMENT {
		err := s.orderRepository.Cancel(ctx, orderUUID)
		if err != nil {
			cancelOrderSpan.RecordError(err)
			if errors.Is(err, repoModel.ErrOrderNotFound) {
				return model.ErrOrderNotFound
			}
			return err
		}
	}

	cancelOrderSpan.SetAttributes(
		attribute.String("order.status", string(order.Status)),
	)

	return nil
}
