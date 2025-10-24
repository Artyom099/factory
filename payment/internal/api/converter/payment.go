package converter

import (
	"github.com/Artyom099/factory/payment/internal/service/model"
	paymentV1 "github.com/Artyom099/factory/shared/pkg/proto/payment/v1"
)

func ToModelPayment(req *paymentV1.PayOrderRequest) model.Payment {
	return model.Payment{
		OrderUuid:     req.GetOrderUuid(),
		UserUuid:      req.GetUserUuid(),
		PaymentMethod: model.PaymentMethod(req.GetPaymentMethod()),
	}
}
