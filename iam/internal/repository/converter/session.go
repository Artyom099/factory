package converter

import (
	"time"

	"github.com/samber/lo"

	"github.com/Artyom099/factory/iam/internal/model"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func ToModelSession(dto repoModel.SessionRedisView) model.Session {
	createdAt := time.Unix(0, dto.CreatedAtNs)

	var updatedAt *time.Time
	if dto.UpdatedAtNs != nil {
		updatedAt = lo.ToPtr(time.Unix(0, *dto.UpdatedAtNs))
	}

	expiredAt := time.Unix(0, dto.ExpiredAtNs)

	return model.Session{
		ID:        dto.UUID,
		UserID:    dto.UserID,
		Login:     dto.Login,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiredAt: expiredAt,
	}
}

func ToRedisViewSession(dto model.Session) repoModel.SessionRedisView {
	createdNs := dto.CreatedAt.UnixNano()

	var updatedNs *int64
	if dto.UpdatedAt != nil {
		updatedNs = lo.ToPtr(dto.UpdatedAt.UnixNano())
	}

	expiredNs := dto.ExpiredAt.UnixNano()

	return repoModel.SessionRedisView{
		UUID:        dto.ID,
		UserID:      dto.UserID,
		Login:       dto.Login,
		CreatedAtNs: createdNs,
		UpdatedAtNs: updatedNs,
		ExpiredAtNs: expiredNs,
	}
}
