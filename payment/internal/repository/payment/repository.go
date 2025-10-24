package payment

import (
	def "github.com/Artyom099/factory/payment/internal/repository"
)

var _ def.IPaymentRepository = (*repository)(nil)

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}
