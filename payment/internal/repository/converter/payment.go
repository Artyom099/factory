package converter

import (
	repoModel "github.com/Artyom099/factory/payment/internal/repository/model"
	servModel "github.com/Artyom099/factory/payment/internal/service/model"
)

func ToRepoPayment(dto servModel.Payment) repoModel.RepoPayment {
	return repoModel.RepoPayment{
		OrderUuid:     dto.OrderUuid,
		UserUuid:      dto.UserUuid,
		PaymentMethod: repoModel.PaymentMethod(dto.PaymentMethod),
	}
}
