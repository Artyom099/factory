package payment

import (
	"context"

	"github.com/Artyom099/factory/payment/internal/repository/converter"
	"github.com/Artyom099/factory/payment/internal/service/model"
)

func (s *service) PayOrder(ctx context.Context, dto model.PayOrderServiceRequestDto) (string, error) {
	transactionUuid, err := s.paymentRepository.PayOrder(ctx, converter.PaymentServiceRequestDtoToPaymentRepoRequestDto(dto))
	if err != nil {
		return "", err
	}

	return transactionUuid, nil
}
