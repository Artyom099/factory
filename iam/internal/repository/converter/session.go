package converter

import (
	"time"

	"github.com/Artyom099/factory/iam/internal/model"
	repoModel "github.com/Artyom099/factory/iam/internal/repository/model"
)

func SessionFromRedisView(dto repoModel.SessionRedisView) model.Session {
	createdAt := time.Unix(0, dto.CreatedAtNs)

	var updatedAt *time.Time
	if dto.UpdatedAtNs != nil {
		t := time.Unix(0, *dto.UpdatedAtNs)
		updatedAt = &t
	}

	expiredAt := time.Unix(0, dto.ExpiredAtNs)

	return model.Session{
		ID:        dto.UUID,
		UserID:    dto.UserID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		ExpiredAt: expiredAt,
	}
}

func SessionToRedisView(dto model.Session) repoModel.SessionRedisView {
	createdNs := dto.CreatedAt.UnixNano()

	var updatedNs *int64
	if dto.UpdatedAt != nil {
		v := dto.UpdatedAt.UnixNano()
		updatedNs = &v
	}

	expiredNs := dto.ExpiredAt.UnixNano()

	return repoModel.SessionRedisView{
		UUID:        dto.ID,
		UserID:      dto.UserID,
		CreatedAtNs: createdNs,
		UpdatedAtNs: updatedNs,
		ExpiredAtNs: expiredNs,
	}
}
