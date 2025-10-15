package payment

import (
	"github.com/Artyom099/factory/payment/internal/repository"
	def "github.com/Artyom099/factory/payment/internal/service"
)

var _ def.IPaymentService = (*service)(nil)

type service struct {
	paymentRepository repository.IPaymentRepository
}

func NewService(paymentRepository repository.IPaymentRepository) *service {
	return &service{
		paymentRepository: paymentRepository,
	}
}
