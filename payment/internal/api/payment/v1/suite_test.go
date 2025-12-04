package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/payment/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	paymentService *mocks.IPaymentService

	api *api
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()

	a.paymentService = mocks.NewIPaymentService(a.T())

	a.api = NewAPI(
		a.paymentService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestUnitPaymentApi(t *testing.T) {
	suite.Run(t, new(APISuite))
}
