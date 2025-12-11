package user

import (
	"errors"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *ServiceSuite) TestRegisterSuccess() {
	var (
		email    = gofakeit.Email()
		password = "StrongPassword123"
		userID   = gofakeit.UUID()

		inputUser = model.User{
			Email:    email,
			Password: password,
		}
	)

	s.userRepository.
		On("Get", s.ctx, email).
		Return(model.User{}, model.ErrUserNotFound).
		Once()

	s.userRepository.
		On("Create", s.ctx, mock.AnythingOfType("model.User")).
		Return(userID, nil).
		Once()

	id, err := s.service.Register(s.ctx, inputUser)

	s.Require().NoError(err)
	s.Require().Equal(userID, id)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterUserAlreadyExists() {
	email := gofakeit.Email()

	existing := model.User{
		ID:    gofakeit.UUID(),
		Email: email,
	}

	s.userRepository.
		On("Get", s.ctx, email).
		Return(existing, nil).
		Once()

	_, err := s.service.Register(s.ctx, model.User{Email: email})

	s.Require().Error(err)
	// Текущая логика → не возвращает отдельную ошибку, но Create не вызывается
	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterGetUnexpectedError() {
	var (
		email   = gofakeit.Email()
		getErr  = errors.New("database down")
		payload = model.User{Email: email}
	)

	s.userRepository.
		On("Get", s.ctx, email).
		Return(model.User{}, getErr).
		Once()

	_, err := s.service.Register(s.ctx, payload)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to check existing user")

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestRegisterCreateError() {
	var (
		email     = gofakeit.Email()
		password  = "12345pass"
		createErr = errors.New("failed to insert")
		payload   = model.User{
			Email:    email,
			Password: password,
		}
	)

	s.userRepository.
		On("Get", s.ctx, email).
		Return(model.User{}, model.ErrUserNotFound).
		Once()

	s.userRepository.
		On("Create", s.ctx, mock.AnythingOfType("model.User")).
		Return("", createErr).
		Once()

	_, err := s.service.Register(s.ctx, payload)

	s.Require().Error(err)
	s.Require().Contains(err.Error(), "failed to create user")

	s.userRepository.AssertExpectations(s.T())
}
