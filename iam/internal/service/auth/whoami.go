package auth

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/iam/internal/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) Whoami(ctx context.Context, sessionUuid string) (model.User, model.Session, error) {
	ctx, getUserSpan := tracing.StartSpan(ctx, "iam.get_user",
		trace.WithAttributes(
			attribute.String("user.sessionUuid", sessionUuid),
		),
	)

	session, err := s.sessionRepository.Get(ctx, sessionUuid)
	if err != nil {
		getUserSpan.RecordError(err)
		getUserSpan.End()
		return model.User{}, model.Session{}, err
	}

	getUserSpan.End()

	ctx, getSessionSpan := tracing.StartSpan(ctx, "iam.get_session",
		trace.WithAttributes(
			attribute.String("user.sessionUuid", sessionUuid),
		),
	)
	defer getSessionSpan.End()

	user, err := s.userRepository.Get(ctx, session.Login)
	if err != nil {
		getSessionSpan.RecordError(err)
		return model.User{}, model.Session{}, err
	}

	return user, session, nil
}
