package converter

import (
	servModel "github.com/Artyom099/factory/payment/internal/model"
	repoModel "github.com/Artyom099/factory/payment/internal/repository/model"
)

func PaymentServiceRequestDtoToPaymentRepoRequestDto(dto servModel.PayOrderServiceRequestDto) repoModel.PayOrderRepoRequestDto {
	return repoModel.PayOrderRepoRequestDto{
		OrderUuid:     dto.OrderUuid,
		UserUuid:      dto.UserUuid,
		PaymentMethod: repoModel.PaymentMethod(dto.PaymentMethod),
	}
}
