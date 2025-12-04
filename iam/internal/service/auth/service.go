package auth

import (
	"time"

	"github.com/Artyom099/factory/iam/internal/repository"
	def "github.com/Artyom099/factory/iam/internal/service"
)

var _ def.IAuthService = (*service)(nil)

type service struct {
	sessionTTL        *time.Duration
	userRepository    repository.IUserRepository
	sessionRepository repository.ISessionRepository
}

func NewService(
	sessionTTL *time.Duration,
	userRepository repository.IUserRepository,
	sessionRepository repository.ISessionRepository,
) *service {
	return &service{
		sessionTTL:        sessionTTL,
		userRepository:    userRepository,
		sessionRepository: sessionRepository,
	}
}
