package order

import (
	"context"

	"github.com/google/uuid"

	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	"github.com/Artyom099/factory/order/internal/service/model"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func (s *service) Pay(ctx context.Context, dto model.OrderPayServiceRequestDto) (model.OrderPayServiceResponseDto, error) {
	order, err := s.orderRepository.Get(ctx, dto.OrderUUID)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderNotFound
	}

	if model.OrderStatus(order.Status) == model.OrderStatusPAID {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderAlreadyPaid
	}

	if model.OrderStatus(order.Status) == model.OrderStatusCANCELLED {
		return model.OrderPayServiceResponseDto{}, model.ErrOrderCancelled
	}

	userUuid := uuid.New().String()
	payOrderDto := paymentV1.PayOrderRequest{
		UserUuid:      userUuid,
		OrderUuid:     dto.OrderUUID,
		PaymentMethod: convertPaymentMethodFromOrderToPayment(dto.PaymentMethod),
	}

	res, err := s.paymentClient.PayOrder(ctx, &payOrderDto)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrInPaymeentService
	}

	updateOrderDto := repoModel.OrderUpdateRepoRequestDto{
		OrderUUID:       dto.OrderUUID,
		Status:          repoModel.OrderStatusPAID,
		TransactionUUID: res.TransactionUuid,
		PaymentMethod:   repoModel.OrderPaymentMethod(dto.PaymentMethod),
	}
	_, err = s.orderRepository.Update(ctx, updateOrderDto)
	if err != nil {
		return model.OrderPayServiceResponseDto{}, model.ErrUpdateOrder
	}

	return model.OrderPayServiceResponseDto{
		TransactionUUID: res.TransactionUuid,
	}, nil
}

func convertPaymentMethodFromOrderToPayment(method model.OrderPaymentMethod) paymentV1.PaymentMethod {
	switch method {
	case model.OrderPaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.OrderPaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.OrderPaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.OrderPaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
