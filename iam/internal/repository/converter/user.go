package converter

import (
	"github.com/Artyom099/factory/iam/internal/model"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func ToModelUser(u repoModel.RepoUser) model.User {
	methods := make([]model.NotificationMethod, 0, len(u.NotificationMethods))

	for _, m := range u.NotificationMethods {
		methods = append(methods, model.NotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}

	return model.User{
		ID:                  u.ID,
		Login:               u.Login,
		Email:               u.Email,
		NotificationMethods: methods,
		CreatedAt:           u.CreatedAt,
		UpdatedAt:           u.UpdatedAt,
	}
}

func ToRepoUser(u model.User) repoModel.RepoUser {
	methods := make([]repoModel.RepoNotificationMethod, 0, len(u.NotificationMethods))
	for _, m := range u.NotificationMethods {
		methods = append(methods, repoModel.RepoNotificationMethod{
			ProviderName: m.ProviderName,
			Target:       m.Target,
		})
	}

	return repoModel.RepoUser{
		ID:                  u.ID,
		Login:               u.Login,
		Email:               u.Email,
		Hash:                u.Hash,
		NotificationMethods: methods,
	}
}
