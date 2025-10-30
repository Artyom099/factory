package payment

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/service/model"
)

func (s *service) PayOrder(ctx context.Context, dto model.Payment) (string, error) {
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
		return "", err
	}

	return transactionUuid, nil
}
