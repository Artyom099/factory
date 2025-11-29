package user

import (
	"github.com/Artyom099/factory/iam/internal/repository"
	def "github.com/Artyom099/factory/iam/internal/service"
)

var _ def.IUserService = (*service)(nil)

type service struct {
	userRepository repository.IUserRepository
}

func NewService(userRepository repository.IUserRepository) *service {
	return &service{
		userRepository: userRepository,
	}
}
