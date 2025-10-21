package order

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderUuid = gofakeit.UUID()

		getRepoRequestDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}

		transactionUUID = gofakeit.UUID()

		updateRepoRequestDto = repoModel.OrderUpdateRepoRequestDto{
			OrderUUID:       orderUuid,
			Status:          repoModel.OrderStatusPAID,
			TransactionUUID: transactionUUID,
			PaymentMethod:   repoModel.OrderPaymentMethodSBP,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, serviceRequestDto.PaymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, nil)
	s.orderRepository.On("Update", s.ctx, updateRepoRequestDto).Return(repoModel.OrderUpdateRepoResponseDto{}, nil)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(res, transactionUUID)
}

func (s *ServiceSuite) TestPayRepoUpdateError() {
	var (
		repoErr   = gofakeit.Error()
		orderUuid = gofakeit.UUID()

		getRepoRequestDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}

		transactionUUID = gofakeit.UUID()

		updateRepoRequestDto = repoModel.OrderUpdateRepoRequestDto{
			OrderUUID:       orderUuid,
			Status:          repoModel.OrderStatusPAID,
			TransactionUUID: transactionUUID,
			PaymentMethod:   repoModel.OrderPaymentMethodSBP,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, serviceRequestDto.PaymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, nil)
	s.orderRepository.On("Update", s.ctx, updateRepoRequestDto).Return(repoModel.OrderUpdateRepoResponseDto{}, repoErr)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrUpdateOrder)
}

func (s *ServiceSuite) TestPayRepoGetError() {
	var (
		repoErr   = gofakeit.Error()
		orderUuid = gofakeit.UUID()

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(repoModel.OrderGetRepoResponseDto{}, repoErr)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderNotFound)
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	var (
		repoErr         = gofakeit.Error()
		orderUuid       = gofakeit.UUID()
		transactionUUID = gofakeit.UUID()

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}

		getRepoRequestDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)
	s.paymentClient.On("PayOrder", s.ctx, serviceRequestDto.PaymentMethod, orderUuid, getRepoRequestDto.UserUUID).Return(transactionUUID, repoErr)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrInPaymeentService)
}

func (s *ServiceSuite) TestPayInvalidStatusPaidError() {
	var (
		orderUuid = gofakeit.UUID()

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}

		getRepoRequestDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusPAID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
}

func (s *ServiceSuite) TestPayInvalidStatusCancelledError() {
	var (
		orderUuid = gofakeit.UUID()

		serviceRequestDto = model.OrderPayServiceRequestDto{
			OrderUUID:     orderUuid,
			PaymentMethod: model.OrderPaymentMethodSBP,
		}

		getRepoRequestDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUuid,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      gofakeit.Float32(),
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodUNSPECIFIED,
			Status:          repoModel.OrderStatusCANCELLED,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUuid).Return(getRepoRequestDto, nil)

	res, err := s.service.Pay(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, model.ErrOrderCancelled)
}
