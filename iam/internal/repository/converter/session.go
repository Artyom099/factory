package converter

import (
	"github.com/Artyom099/factory/iam/internal/model"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func SessionFromRedisView(dto repoModel.SessionRedisView) model.Session {
	return model.Session{
		ID: dto.UUID,
	}
}

func SessionToRedisView(dto model.Session) repoModel.SessionRedisView {
	return repoModel.SessionRedisView{
		UUID: dto.ID,
	}
}
