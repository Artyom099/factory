package part

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/inventory/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Create(ctx context.Context, dto model.Part) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "part.create",
		trace.WithAttributes(
			attribute.String("part.uuid", dto.Uuid),
		),
	)
	defer span.End()

	partUUID, err := s.partRepository.Create(ctx, dto)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	return partUUID, nil
}
