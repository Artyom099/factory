package payment

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/payment/internal/repository/model"
	"github.com/Artyom099/factory/payment/internal/service/model"
)

func (s *ServiceSuite) TestPayOrderSuccess() {
	var (
		transactionUUID   = gofakeit.UUID()
		orderUuid         = gofakeit.UUID()
		userUuid          = gofakeit.UUID()
		paymentMethod     = model.PaymentMethod(gofakeit.Number(1, 4))
		repoPaymentMethod = repoModel.PaymentMethod(paymentMethod)

		serviceRequestDto = model.PayOrderServiceRequestDto{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}

		repoRequestDto = repoModel.PayOrderRepoRequestDto{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: repoPaymentMethod,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, repoRequestDto).Return(transactionUUID, nil)

	uuid, err := s.service.PayOrder(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, uuid)
}

func (s *ServiceSuite) TestPayOrderRepoError() {
	var (
		repoErr           = gofakeit.Error()
		orderUuid         = gofakeit.UUID()
		userUuid          = gofakeit.UUID()
		paymentMethod     = model.PaymentMethod(gofakeit.Number(1, 4))
		repoPaymentMethod = repoModel.PaymentMethod(paymentMethod)

		serviceRequestDto = model.PayOrderServiceRequestDto{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: paymentMethod,
		}

		repoRequestDto = repoModel.PayOrderRepoRequestDto{
			OrderUuid:     orderUuid,
			UserUuid:      userUuid,
			PaymentMethod: repoPaymentMethod,
		}
	)

	s.paymentRepository.On("PayOrder", s.ctx, repoRequestDto).Return("", repoErr)

	uuid, err := s.service.PayOrder(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(uuid)
}
