package part

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Get(ctx context.Context, uuid string) (servModel.Part, error) {
	ctx, span := tracing.StartSpan(ctx, "part.get",
		trace.WithAttributes(
			attribute.String("part.uuid", uuid),
		),
	)
	defer span.End()

	res, err := s.partRepository.Get(ctx, uuid)
	if err != nil {
		span.RecordError(err)
		if errors.Is(err, repoModel.ErrPartNotFound) {
			return servModel.Part{}, servModel.ErrPartNotFound
		}
		return servModel.Part{}, err
	}

	return res, nil
}
