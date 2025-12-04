package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Artyom099/factory/iam/internal/model"
	commonV1 "github.com/Artyom099/factory/shared/pkg/proto/common/v1"
	userV1 "github.com/Artyom099/factory/shared/pkg/proto/user/v1"
)

func ToApiUser(user model.User) *userV1.GetUserResponse {
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
		apiUser.Info.NotificationMethods = append(
			apiUser.Info.NotificationMethods,
			&commonV1.NotificationMethod{
				ProviderName: nm.ProviderName,
				Target:       nm.Target,
			},
		)
	}

	return &userV1.GetUserResponse{
		User: apiUser,
	}
}

func ToModelRegUser(req *userV1.RegisterRequest) model.User {
	info := req.GetInfo()
	userInfo := info.GetInfo()

	user := model.User{
		Login:    userInfo.GetLogin(),
		Email:    userInfo.GetEmail(),
		Password: info.GetPassword(),
		// Hash: will be set later in service layer after hashing
	}

	for _, nm := range userInfo.GetNotificationMethods() {
		user.NotificationMethods = append(user.NotificationMethods, model.NotificationMethod{
			ProviderName: nm.GetProviderName(),
			Target:       nm.GetTarget(),
		})
	}

	return user
}
