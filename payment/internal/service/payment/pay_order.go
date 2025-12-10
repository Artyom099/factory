package payment

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/Artyom099/factory/payment/internal/service/model"
	"github.com/Artyom099/factory/platform/pkg/tracing"
)

func (s *service) PayOrder(ctx context.Context, dto model.Payment) (string, error) {
	ctx, span := tracing.StartSpan(ctx, "payment.pay_order",
		trace.WithAttributes(
			attribute.String("order.uuid", dto.OrderUuid),
			attribute.String("user.uuid", dto.UserUuid),
		),
	)
	defer span.End()

	switch dto.PaymentMethod {
	case model.CARD,
		model.SBP,
		model.CREDIT_CARD,
		model.INVESTOR_MONEY:
	default:
		return "", model.ErrInvalidPaymentMethod
	}

	transactionUuid, err := s.paymentRepository.PayOrder(ctx, dto)
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	span.SetAttributes(
		attribute.String("transactionUuid", transactionUuid),
	)

	return transactionUuid, nil
}
