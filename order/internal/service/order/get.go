package order

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Get(ctx context.Context, orderUUID string) (servModel.Order, error) {
	ctx, span := tracing.StartSpan(ctx, "order.get",
		trace.WithAttributes(
			attribute.String("order.uuid", orderUUID),
		),
	)
	defer span.End()

	res, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, repoModel.ErrOrderNotFound) {
			return servModel.Order{}, servModel.ErrOrderNotFound
		}
		return servModel.Order{}, err
	}

	return res, nil
}
