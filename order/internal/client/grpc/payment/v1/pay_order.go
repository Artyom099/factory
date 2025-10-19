package v1

import (
	"context"

	"github.com/Artyom099/factory/order/internal/service/model"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(
	ctx context.Context,
	paymentMethod model.OrderPaymentMethod,
	orderUUID, userUUID string,
) (string, error) {
	res, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: convertPaymentMethodFromOrderToPayment(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
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
