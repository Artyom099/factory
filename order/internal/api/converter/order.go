package converter

import (
	servModel "github.com/Artyom099/factory/order/internal/service/model"
	apiModel "github.com/Artyom099/factory/shared/pkg/openapi/order/v1"
)

// api - service

func OrderGetApiRequestDtoToOrderGetServiceRequestDto(dto apiModel.GetOrderParams) servModel.OrderGetServiceRequestDto {
	return servModel.OrderGetServiceRequestDto{
		OrderUUID: dto.OrderUUID.String(),
	}
}

func OrderGetServiceResponseDtoToOrderGetApiResponseDto(dto servModel.OrderGetServiceResponseDto) *apiModel.Order {
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

func OrderCreateApiRequestDtoToOrderCreateServiceRequestDto(dto *apiModel.OrderCreateRequest) servModel.OrderCreateServiceRequestDto {
	return servModel.OrderCreateServiceRequestDto{
		UserUUID:  dto.UserUUID,
		PartUuids: dto.PartUuids,
	}
}

func OrderCancelApiRequestDtoToOrderCancelServiceRequestDto(dto apiModel.CancelOrderParams) string {
	return dto.OrderUUID.String()
}

func OrderPayApiRequestDtoToOrderPayServiceRequestDto(
	param apiModel.PayOrderParams,
	req *apiModel.OrderPayRequest,
) servModel.OrderPayServiceRequestDto {
	return servModel.OrderPayServiceRequestDto{
		OrderUUID:     param.OrderUUID.String(),
		PaymentMethod: servModel.OrderPaymentMethod(req.PaymentMethod),
	}
}
