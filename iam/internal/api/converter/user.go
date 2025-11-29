package converter

import (
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
		// CreatedAt: &timestamppb.Timestamp{},
	}

	// if user.UpdatedAt != nil {
	// apiUser.UpdatedAt =
	// }

	return &userV1.GetUserResponse{
		User: apiUser,
	}
}
