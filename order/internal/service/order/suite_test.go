package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/Artyom099/factory/order/internal/client/grpc/mocks"
	"github.com/Artyom099/factory/order/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	orderRepository *mocks.IOrderRepository
	inventoryClient *clientMocks.IInventoryClient
	paymentClient   *clientMocks.IPaymentClient

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = mocks.NewIOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewIInventoryClient(s.T())
	s.paymentClient = clientMocks.NewIPaymentClient(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
