package order

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	clientMocks "github.com/Artyom099/factory/order/internal/client/grpc/mocks"
	repoMocks "github.com/Artyom099/factory/order/internal/repository/mocks"
	serviceMocks "github.com/Artyom099/factory/order/internal/service/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	orderRepository      *repoMocks.IOrderRepository
	inventoryClient      *clientMocks.IInventoryClient
	paymentClient        *clientMocks.IPaymentClient
	orderProducerService *serviceMocks.IOrderProducerService
	orderConsumerService *serviceMocks.IOrderConsumerService

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewIOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewIInventoryClient(s.T())
	s.paymentClient = clientMocks.NewIPaymentClient(s.T())
	s.orderProducerService = serviceMocks.NewIOrderProducerService(s.T())
	s.orderConsumerService = serviceMocks.NewIOrderConsumerService(s.T())

	s.service = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
		s.orderProducerService,
		s.orderConsumerService,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
