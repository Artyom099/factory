package converter

import (
	repoModel "github.com/Artyom099/factory/payment/internal/repository/model"
	servModel "github.com/Artyom099/factory/payment/internal/service/model"
)

func PaymentServiceRequestDtoToPaymentRepoRequestDto(dto servModel.PayOrderServiceRequestDto) repoModel.PayOrderRepoRequestDto {
	return repoModel.PayOrderRepoRequestDto{
		OrderUuid:     dto.OrderUuid,
		UserUuid:      dto.UserUuid,
		PaymentMethod: repoModel.PaymentMethod(dto.PaymentMethod),
	}
}
