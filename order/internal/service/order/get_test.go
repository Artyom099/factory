package order

import (
	"github.com/brianvoe/gofakeit/v6"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		serviceResponseDto = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      100.0,
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodCARD,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(serviceResponseDto, nil)

	res, err := s.service.Get(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(res, serviceResponseDto)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr   = gofakeit.Error()
		orderUUID = gofakeit.UUID()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(model.Order{}, repoErr)

	res, err := s.service.Get(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}
