package payment

import (
	def "github.com/Artyom099/factory/payment/internal/repository"
	repoModel "github.com/Artyom099/factory/payment/internal/repository/model"
)

var _ def.IPaymentRepository = (*repository)(nil)

type repository struct {
	data map[string]repoModel.PayOrderRepoRequestDto // todo мапа здесь не нужна
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.PayOrderRepoRequestDto),
	}
}
