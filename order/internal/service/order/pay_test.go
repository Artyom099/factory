package order

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"

	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUuid     = gofakeit.UUID()
		paymentMethod = model.OrderPaymentMethodSBP

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}

		transactionUUID = gofakeit.UUID()

		updateRepoRequestDto = getRepoRequestDto
	)
	updateRepoRequestDto.Status = model.OrderStatusPAID
	updateRepoRequestDto.TransactionUUID = transactionUUID
	updateRepoRequestDto.PaymentMethod = model.OrderPaymentMethodSBP

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, paymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, nil)
	s.orderRepository.On("Update", s.ctx, updateRepoRequestDto).Return(nil)
	s.orderProducerService.On("ProduceOrderPaid", s.ctx, mock.AnythingOfType("model.OrderPaidOutEvent")).Return(nil)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().NoError(err)
	s.Require().Equal(res, transactionUUID)
}

func (s *ServiceSuite) TestPayRepoUpdateError() {
	var (
		repoErr       = gofakeit.Error()
		orderUuid     = gofakeit.UUID()
		paymentMethod = model.OrderPaymentMethodSBP

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}

		transactionUUID = gofakeit.UUID()

		updateRepoRequestDto = getRepoRequestDto
	)
	updateRepoRequestDto.Status = model.OrderStatusPAID
	updateRepoRequestDto.TransactionUUID = transactionUUID
	updateRepoRequestDto.PaymentMethod = model.OrderPaymentMethodSBP

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, paymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, nil)
	s.orderRepository.On("Update", s.ctx, updateRepoRequestDto).Return(repoErr)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayRepoGetError() {
	var (
		repoErr       = gofakeit.Error()
		orderUuid     = gofakeit.UUID()
		paymentMethod = model.OrderPaymentMethodSBP
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(model.Order{}, repoErr)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	var (
		repoErr         = gofakeit.Error()
		orderUuid       = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()
		paymentMethod   = model.OrderPaymentMethodSBP

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, paymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, repoErr)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrInPaymeentService)
}

func (s *ServiceSuite) TestPayInvalidStatusPaidError() {
	var (
		orderUuid     = gofakeit.UUID()
		paymentMethod = model.OrderPaymentMethodSBP

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayInvalidStatusCancelledError() {
	var (
		orderUuid     = gofakeit.UUID()
		paymentMethod = model.OrderPaymentMethodSBP

		getRepoRequestDto = model.Order{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   model.OrderPaymentMethodUNSPECIFIED,
			Status:          model.OrderStatusCANCELLED,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)

	res, err := s.service.Pay(s.ctx, orderUuid, paymentMethod)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}
