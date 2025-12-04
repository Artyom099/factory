package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Artyom099/factory/iam/internal/model"
	authV1 "github.com/Artyom099/factory/shared/pkg/proto/auth/v1"
	commonV1 "github.com/Artyom099/factory/shared/pkg/proto/common/v1"
)

func ToApiLogin(sessionUuid string) *authV1.LoginResponse {
	return &authV1.LoginResponse{
		SessionUuid: sessionUuid,
	}
}

func ToApiWhoami(user model.User, session model.Session) *authV1.WhoamiResponse {
	apiUser := &commonV1.User{
		Uuid: user.ID,
		Info: &commonV1.UserInfo{
			Login: user.Login,
			Email: user.Email,
		},
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		apiUser.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	for _, nm := range user.NotificationMethods {
		apiUser.Info.NotificationMethods = append(apiUser.Info.NotificationMethods, &commonV1.NotificationMethod{
			ProviderName: nm.ProviderName,
			Target:       nm.Target,
		})
	}

	apiSession := &commonV1.Session{
		Uuid:      session.ID,
		CreatedAt: timestamppb.New(session.CreatedAt),
		ExpiresAt: timestamppb.New(session.ExpiredAt),
	}

	if session.UpdatedAt != nil {
		apiSession.UpdatedAt = timestamppb.New(*session.UpdatedAt)
	}

	return &authV1.WhoamiResponse{
		User:    apiUser,
		Session: apiSession,
	}
}
