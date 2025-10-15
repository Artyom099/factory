package part

import (
	"context"
	"time"

	"github.com/google/uuid"

	repoModel "github.com/Artyom099/factory/inventory/internal/repository/model"
)

func (r *repository) Create(ctx context.Context, dto repoModel.PartCreateRepoRequest) (string, error) {
	uuid := uuid.New().String()
	dto.Part.Uuid = uuid
	dto.Part.CreatedAt = time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[uuid] = dto.Part

	return uuid, nil
}
