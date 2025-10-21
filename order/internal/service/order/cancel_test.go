package order

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(getRepoResponseDto, nil)
	s.orderRepository.On("Cancel", s.ctx, orderUUID).Return(repoModel.OrderCancelRepoResponseDto{}, nil)

	res, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
	s.Require().Equal(res, model.OrderCancelServiceResponseDto{})
}

func (s *ServiceSuite) TestCancelRepoCancelError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()

		getRepoResponseDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(getRepoResponseDto, nil)
	s.orderRepository.On("Cancel", s.ctx, orderUUID).Return(repoModel.OrderCancelRepoResponseDto{}, repoErr)

	res, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancelConflictError() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(getRepoResponseDto, nil)

	res, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrConflict)
}

func (s *ServiceSuite) TestCancelRepoGetError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(repoModel.OrderGetRepoResponseDto{}, repoErr)

	res, err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}
