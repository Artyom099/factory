package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/payment/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	paymentRepository *mocks.IPaymentRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.paymentRepository = mocks.NewIPaymentRepository(s.T())

	s.service = NewService(
		s.paymentRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
