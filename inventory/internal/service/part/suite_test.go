package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	partRepository *mocks.IPartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.partRepository = mocks.NewIPartRepository(s.T())

	s.service = NewService(
		s.partRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestUnitPartService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
