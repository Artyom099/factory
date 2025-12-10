package part

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/inventory/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) List(ctx context.Context, dto model.PartFilter) ([]model.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "inventory.list_parts",
		trace.WithAttributes(
			attribute.StringSlice("part.uuids", dto.Uuids),
		),
	)
	defer span.End()

	res, err := s.partRepository.List(ctx, dto)
	if err != nil {
		span.RecordError(err)
		return []model.Part{}, err
	}

	span.SetAttributes(
		attribute.Int("parts.count", len(res)),
	)

	return res, nil
}
