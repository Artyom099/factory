package payment

import (
	"github.com/brianvoe/gofakeit/v6"

	"github.com/Artyom099/factory/payment/internal/service/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		transactionUUID = gofakeit.UUID()
		orderUuid       = gofakeit.UUID()
		userUuid        = gofakeit.UUID()
		paymentMethod   = model.PaymentMethod(gofakeit.Number(1, 4))

		serviceRequestDto = model.Payment{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, serviceRequestDto).Return(transactionUUID, nil)

	uuid, err := s.service.PayOrder(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, uuid)
}

func (s *ServiceSuite) TestPayOrderRepoError() {
	var (
		repoErr       = gofakeit.Error()
		orderUuid     = gofakeit.UUID()
		userUuid      = gofakeit.UUID()
		paymentMethod = model.PaymentMethod(gofakeit.Number(1, 4))

		serviceRequestDto = model.Payment{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, serviceRequestDto).Return("", repoErr)

	uuid, err := s.service.PayOrder(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrInternalError)
	s.Require().Empty(uuid)
}
