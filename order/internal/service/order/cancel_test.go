package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", mock.Anything, orderUUID).Return(getRepoResponseDto, nil)
	s.orderRepository.On("Cancel", mock.Anything, orderUUID).Return(nil)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelRepoCancelError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()

		getRepoResponseDto = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", mock.Anything, orderUUID).Return(getRepoResponseDto, nil)
	s.orderRepository.On("Cancel", mock.Anything, orderUUID).Return(repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancelConflictError() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", mock.Anything, orderUUID).Return(getRepoResponseDto, nil)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrConflict)
}

func (s *ServiceSuite) TestCancelRepoGetError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("Get", mock.Anything, orderUUID).Return(model.Order{}, repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
