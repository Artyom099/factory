package v1

import (
	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Artyom099/factory/payment/internal/api/converter"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func (s *APISuite) TestPayOrderSuccess() {
	var (
		transactionUUID = gofakeit.UUID()

		apiRequestDto = &paymentV1.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod(gofakeit.Number(1, 4)),
		}

		expectedModelInfo = converter.ToModelPayment(apiRequestDto)
	)

	s.paymentService.On("PayOrder", s.ctx, expectedModelInfo).Return(transactionUUID, nil)

	res, err := s.api.PayOrder(s.ctx, apiRequestDto)
	s.Require().NoError(err)
	s.Require().NotNil(res)
}

func (s *APISuite) TestPayOrderServiceError() {
	var (
		serviceErr = gofakeit.Error()

		apiRequestDto = &paymentV1.PayOrderRequest{
			OrderUuid:     gofakeit.UUID(),
			UserUuid:      gofakeit.UUID(),
			PaymentMethod: paymentV1.PaymentMethod(gofakeit.Number(1, 4)),
		}

		expectedModelInfo = converter.ToModelPayment(apiRequestDto)
	)

	s.paymentService.On("PayOrder", s.ctx, expectedModelInfo).Return("", serviceErr)

	res, err := s.api.PayOrder(s.ctx, apiRequestDto)
	s.Require().Error(err)
	s.Require().ErrorIs(err, status.Error(codes.Internal, "internal server error"))
	s.Require().Nil(res)
}
