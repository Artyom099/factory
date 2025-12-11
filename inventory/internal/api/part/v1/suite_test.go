package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/inventory/internal/service/mocks"
)

type APISuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	partService *mocks.IPartService

	api *api
}

func (a *APISuite) SetupTest() {
	a.ctx = context.Background()

	a.partService = mocks.NewIPartService(a.T())

	a.api = NewAPI(
		a.partService,
	)
}

func (s *APISuite) TearDownTest() {
}

func TestUnitPartApi(t *testing.T) {
	suite.Run(t, new(APISuite))
}
