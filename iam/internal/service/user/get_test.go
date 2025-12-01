package user

import (
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v6"

	"github.com/Artyom099/factory/iam/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		userUUID = gofakeit.UUID()

		notificationMethods = []model.NotificationMethod{
			{ProviderName: "telegram", Target: "@test_get"},
			{ProviderName: "email", Target: "test@gmail.com"},
		}

		modelUser = model.User{
			ID:                  userUUID,
			Login:               gofakeit.Username(),
			Email:               gofakeit.Email(),
			Password:            gofakeit.Password(true, true, true, true, false, 15),
			Hash:                gofakeit.HackerVerb(),
			NotificationMethods: notificationMethods,
			CreatedAt:           time.Now(),
			UpdatedAt:           nil,
		}
	)

	s.userRepository.On("Get", s.ctx, userUUID).Return(modelUser, nil)

	res, err := s.service.Get(s.ctx, userUUID)
	s.Require().NoError(err)

	s.Require().Equal(userUUID, res.ID)
	s.Require().Equal(modelUser.Login, res.Login)
	s.Require().Equal(modelUser.Email, res.Email)
	s.Require().Equal(modelUser.Hash, res.Hash)
	s.Require().Equal(modelUser.Password, res.Password)
	s.Require().Equal(modelUser.CreatedAt, res.CreatedAt)
	s.Require().Nil(res.UpdatedAt, "UpdatedAt должен быть nil")

	s.Require().Len(res.NotificationMethods, len(notificationMethods))
	s.Require().Equal(notificationMethods, res.NotificationMethods)

	s.userRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetError() {
	var (
		userUUID = gofakeit.UUID()
		repoErr  = errors.New("database failure")
	)

	s.userRepository.
		On("Get", s.ctx, userUUID).
		Return(model.User{}, repoErr).
		Once()

	res, err := s.service.Get(s.ctx, userUUID)

	s.Require().Error(err)
	s.Require().EqualError(err, "database failure")

	s.Require().Equal(model.User{}, res)

	s.userRepository.AssertExpectations(s.T())
}
