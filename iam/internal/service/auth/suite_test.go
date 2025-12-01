package auth

import (
	"context"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"

	"github.com/Artyom099/factory/iam/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context // nolint:containedctx

	sessionTTL        *time.Duration
	userRepository    *mocks.IUserRepository
	sessionRepository *mocks.ISessionRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.sessionTTL = lo.ToPtr(time.Duration(10 * time.Minute))
	s.userRepository = mocks.NewIUserRepository(s.T())
	s.sessionRepository = mocks.NewISessionRepository(s.T())

	s.service = NewService(
		s.sessionTTL,
		s.userRepository,
		s.sessionRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestUnitAuthService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
