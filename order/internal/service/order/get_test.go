package order

import (
	"github.com/brianvoe/gofakeit/v6"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderUUID = gofakeit.UUID()

		serviceRequestDto = model.OrderGetServiceRequestDto{
			OrderUUID: orderUUID,
		}

		repoResponseDto = repoModel.OrderGetRepoResponseDto{
			OrderUUID:       orderUUID,
			UserUUID:        gofakeit.UUID(),
			PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
			TotalPrice:      100.0,
			TransactionUUID: gofakeit.UUID(),
			PaymentMethod:   repoModel.OrderPaymentMethodCARD,
			Status:          repoModel.OrderStatusPENDINGPAYMENT,
		}

		serviceResponseDto = model.OrderGetServiceResponseDto{
			OrderUUID:       orderUUID,
			UserUUID:        repoResponseDto.UserUUID,
			PartUuids:       repoResponseDto.PartUuids,
			TotalPrice:      repoResponseDto.TotalPrice,
			TransactionUUID: repoResponseDto.TransactionUUID,
			PaymentMethod:   model.OrderPaymentMethod(repoResponseDto.PaymentMethod),
			Status:          model.OrderStatus(repoResponseDto.Status),
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(repoResponseDto, nil)

	res, err := s.service.Get(s.ctx, serviceRequestDto)
	s.Require().NoError(err)
	s.Require().Equal(res, serviceResponseDto)
}

func (s *ServiceSuite) TestGetRepoError() {
	var (
		repoErr   = gofakeit.Error()
		orderUUID = gofakeit.UUID()

		serviceRequestDto = model.OrderGetServiceRequestDto{
			OrderUUID: orderUUID,
		}
	)

	s.orderRepository.On("Get", s.ctx, orderUUID).Return(repoModel.OrderGetRepoResponseDto{}, repoErr)

	res, err := s.service.Get(s.ctx, serviceRequestDto)
	s.Require().Error(err)
	s.Require().Empty(res)
	s.Require().ErrorIs(err, repoErr)
}
