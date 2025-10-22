package order

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = repoModel.RepoOrder{
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
	s.orderRepository.On("Cancel", s.ctx, orderUUID).Return(nil)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestCancelRepoCancelError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()

		getRepoResponseDto = repoModel.RepoOrder{
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
	s.orderRepository.On("Cancel", s.ctx, orderUUID).Return(repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestCancelConflictError() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = repoModel.RepoOrder{
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

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, model.ErrConflict)
}

func (s *ServiceSuite) TestCancelRepoGetError() {
	var (
		orderUUID = gofakeit.UUID()
		repoErr   = gofakeit.Error()
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(repoModel.RepoOrder{}, repoErr)

	err := s.service.Cancel(s.ctx, orderUUID)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
