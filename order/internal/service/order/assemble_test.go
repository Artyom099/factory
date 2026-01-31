package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestAssembleSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		getRepoResponseDto = model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      float32(gofakeit.Float32()),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodSBP,
			Status:          model.OrderStatusPAID,
		}

		updateRepoRequestDto = getRepoResponseDto
	)
	updateRepoRequestDto.Status = model.OrderStatusASSEMBLED

	s.orderRepository.On("Get", mock.Anything, orderUUID).Return(getRepoResponseDto, nil)
	s.orderRepository.On("Update", mock.Anything, updateRepoRequestDto).Return(nil)

	err := s.service.Assemble(s.ctx, orderUUID)
	s.Require().NoError(err)
}

func (s *ServiceSuite) TestAssembledRepoUpdateError() {
	var (
		repoErr   = gofakeit.Error()
		orderUuid = gofakeit.UUID()

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodSBP,
			Status:          model.OrderStatusPAID,
		}

		updateRepoRequestDto = getRepoRequestDto
	)
	updateRepoRequestDto.Status = model.OrderStatusASSEMBLED

	s.orderRepository.On("Get", mock.Anything, orderUuid).Return(getRepoRequestDto, nil)
	s.orderRepository.On("Update", mock.Anything, updateRepoRequestDto).Return(repoErr)

	err := s.service.Assemble(s.ctx, orderUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestAssembledRepoGetError() {
	var (
		repoErr   = gofakeit.Error()
		orderUuid = gofakeit.UUID()
	)

	s.orderRepository.On("Get", mock.Anything, orderUuid).Return(model.Order{}, repoErr)

	err := s.service.Assemble(s.ctx, orderUuid)
	s.Require().Error(err)
	s.Require().ErrorIs(err, repoErr)
}
