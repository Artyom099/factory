package converter

import (
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
	// todo
	apiUser := &commonV1.User{
		Uuid: user.ID,
		Info: &commonV1.UserInfo{
			Login: user.Login,
			Email: user.Email,
		},
		// CreatedAt: &timestamppb.Timestamp{},
	}

	// if user.UpdatedAt != nil {
	// apiUser.UpdatedAt =
	// }

	// var updatedAt *time.Time
	// if redisView.UpdatedAtNs != nil {
	// 	tmp := time.Unix(0, *redisView.UpdatedAtNs)
	// 	updatedAt = &tmp
	// }

	return &authV1.WhoamiResponse{
		User:    apiUser,
		Session: &commonV1.Session{},
	}
}
