package converter

import (
	repoModel "github.com/Artyom099/factory/order/internal/repository/model"
	servModel "github.com/Artyom099/factory/order/internal/service/model"
)

func OrderCreateServiceRequestDtoToOrderCreateRepoRequestDto(
	dto servModel.OrderCreateServiceRequestDto,
	totalPrice float32,
) repoModel.OrderCreateRepoRequestDto {
	return repoModel.OrderCreateRepoRequestDto{
		UserUUID:   dto.UserUUID,
		PartUuids:  dto.PartUuids,
		TotalPrice: totalPrice,
	}
}

func OrderGetServiceRequestDtoToOrderGetRepoRequestDto(dto servModel.OrderGetServiceRequestDto) repoModel.OrderGetRepoRequestDto {
	return repoModel.OrderGetRepoRequestDto{
		OrderUUID: dto.OrderUUID,
	}
}

func OrderGetRepoResponseDtoToOrderGetServiceResponseDto(dto repoModel.OrderGetRepoResponseDto) servModel.OrderGetServiceResponseDto {
	return servModel.OrderGetServiceResponseDto{
		OrderUUID:       dto.OrderUUID,
		UserUUID:        dto.UserUUID,
		PartUuids:       dto.PartUuids,
		TotalPrice:      dto.TotalPrice,
		TransactionUUID: dto.TransactionUUID,
		PaymentMethod:   servModel.OrderPaymentMethod(dto.PaymentMethod),
		Status:          servModel.OrderStatus(dto.Status),
	}
}
