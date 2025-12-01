package auth

import (
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *ServiceSuite) TestLoginSuccess() {
	var (
		login    = gofakeit.Email()
		password = "MyStrongPass123"

		userID = gofakeit.UUID()
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := model.User{
		ID:    userID,
		Login: login,
		Hash:  string(hash),
	}

	s.userRepository.
		On("Get", s.ctx, login).
		Return(user, nil).
		Once()

	ttl := 10 * time.Minute
	s.service.sessionTTL = &ttl

	s.sessionRepository.
		On("Create", s.ctx, mock.AnythingOfType("model.Session"), ttl).
		Return(nil).
		Once()

	s.sessionRepository.
		On("AddSessionToUserSet", s.ctx, userID, mock.AnythingOfType("string")).
		Return(nil).
		Once()

	sessionID, err := s.service.Login(s.ctx, login, password)

	s.Require().NoError(err)
	s.Require().NotEmpty(sessionID)

	s.userRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLoginUserNotFound() {
	var (
		login    = gofakeit.Email()
		password = "pass"
	)

	s.userRepository.
		On("Get", s.ctx, login).
		Return(model.User{}, model.ErrUserNotFound).
		Once()

	_, err := s.service.Login(s.ctx, login, password)

	s.Require().Error(err)
	s.Require().Equal(model.ErrUserNotFound, err)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLoginInvalidPassword() {
	var (
		login    = gofakeit.Email()
		password = "WrongPassword"
	)

	correctPassword := "Correct123!"
	hash, _ := bcrypt.GenerateFromPassword([]byte(correctPassword), bcrypt.DefaultCost)

	user := model.User{
		ID:    gofakeit.UUID(),
		Login: login,
		Hash:  string(hash),
	}

	s.userRepository.
		On("Get", s.ctx, login).
		Return(user, nil).
		Once()

	_, err := s.service.Login(s.ctx, login, password)

	s.Require().Error(err)
	s.Require().Equal(model.ErrInvalidPassword, err)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLoginSessionCreateError() {
	var (
		login    = gofakeit.Email()
		userID   = gofakeit.UUID()
		password = "MyPass123"
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := model.User{
		ID:    userID,
		Login: login,
		Hash:  string(hash),
	}

	s.userRepository.
		On("Get", s.ctx, login).
		Return(user, nil).
		Once()

	ttl := 5 * time.Minute
	s.service.sessionTTL = &ttl

	createErr := errors.New("redis unavailable")

	s.sessionRepository.
		On("Create", s.ctx, mock.AnythingOfType("model.Session"), ttl).
		Return(createErr).
		Once()

	_, err := s.service.Login(s.ctx, login, password)

	s.Require().Error(err)
	s.Require().Equal(createErr, err)

	s.userRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestLoginAddSessionToUserSetError() {
	var (
		login    = gofakeit.Email()
		userID   = gofakeit.UUID()
		password = "MyPass123"
	)

	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := model.User{
		ID:    userID,
		Login: login,
		Hash:  string(hash),
	}

	s.userRepository.
		On("Get", s.ctx, login).
		Return(user, nil).
		Once()

	ttl := 5 * time.Minute
	s.service.sessionTTL = &ttl

	s.sessionRepository.
		On("Create", s.ctx, mock.AnythingOfType("model.Session"), ttl).
		Return(nil).
		Once()

	addErr := errors.New("redis set error")

	s.sessionRepository.
		On("AddSessionToUserSet", s.ctx, userID, mock.AnythingOfType("string")).
		Return(addErr).
		Once()

	_, err := s.service.Login(s.ctx, login, password)

	s.Require().Error(err)
	s.Require().Equal(addErr, err)

	s.userRepository.AssertExpectations(s.T())
	s.sessionRepository.AssertExpectations(s.T())
}
