package user

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/iam/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	userRepository *mocks.IUserRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.userRepository = mocks.NewIUserRepository(s.T())

	s.service = NewService(
		s.userRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestUnitUserService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
