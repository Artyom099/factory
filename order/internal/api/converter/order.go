package converter

import (
	servModel "github.com/Artyom099/factory/order/internal/service/model"
	apiModel "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

func ModelToApiOrder(dto servModel.Order) *apiModel.Order {
	return &apiModel.Order{
		OrderUUID:       dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: dto.TransactionUUID,
		PaymentMethod:   apiModel.OrderPaymentMethod(dto.PaymentMethod),
		Status:          apiModel.OrderStatus(dto.Status),
	}
}

func ApiToModelOrder(dto *apiModel.OrderCreateRequest) servModel.Order {
	return servModel.Order{
		UserUUID:  dto.UserUUID,
		PartUuids: dto.PartUuids,
	}
}

func ParamsAndReqToModelOrder(
	param apiModel.PayOrderParams,
	req *apiModel.OrderPayRequest,
) servModel.Order {
	return servModel.Order{
		OrderUUID:     param.OrderUUID.String(),
		PaymentMethod: servModel.OrderPaymentMethod(req.PaymentMethod),
	}
}
