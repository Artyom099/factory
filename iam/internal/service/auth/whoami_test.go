package auth

import (
	"errors"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *ServiceSuite) TestWhoamiSuccess() {
	sessionUUID := gofakeit.UUID()

	session := model.Session{
		ID:     sessionUUID,
		Login:  gofakeit.Email(),
		UserID: gofakeit.UUID(),
	}

	user := model.User{
		ID:    session.UserID,
		Login: session.Login,
		Email: gofakeit.Email(),
	}

	s.sessionRepository.
		On("Get", s.ctx, sessionUUID).
		Return(session, nil).
		Once()

	s.userRepository.
		On("Get", s.ctx, session.Login).
		Return(user, nil).
		Once()

	resUser, resSession, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().NoError(err)
	s.Require().Equal(user, resUser)
	s.Require().Equal(session, resSession)

	s.sessionRepository.AssertExpectations(s.T())
	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoamiSessionNotFound() {
	sessionUUID := gofakeit.UUID()

	errSession := errors.New("session not found")

	s.sessionRepository.
		On("Get", s.ctx, sessionUUID).
		Return(model.Session{}, errSession).
		Once()

	_, _, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().Error(err)
	s.Require().Equal(errSession, err)

	s.sessionRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestWhoamiUserNotFound() {
	sessionUUID := gofakeit.UUID()

	session := model.Session{
		ID:     sessionUUID,
		Login:  gofakeit.Email(),
		UserID: gofakeit.UUID(),
	}

	errUser := errors.New("user not found")

	s.sessionRepository.
		On("Get", s.ctx, sessionUUID).
		Return(session, nil).
		Once()

	s.userRepository.
		On("Get", s.ctx, session.Login).
		Return(model.User{}, errUser).
		Once()

	_, _, err := s.service.Whoami(s.ctx, sessionUUID)

	s.Require().Error(err)
	s.Require().Equal(errUser, err)

	s.sessionRepository.AssertExpectations(s.T())
	s.userRepository.AssertExpectations(s.T())
}
