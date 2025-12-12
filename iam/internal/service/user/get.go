package user

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Get(ctx context.Context, uuid string) (model.User, error) {
	ctx, span := tracing.StartSpan(ctx, "iam.get_user",
		trace.WithAttributes(
			attribute.String("user.uuid", uuid),
		),
	)
	defer span.End()

	user, err := s.userRepository.Get(ctx, uuid)
	if err != nil {
		span.RecordError(err)
		return model.User{}, err
	}

	return user, nil
}
