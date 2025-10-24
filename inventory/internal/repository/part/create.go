package part

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	"github.com/Artyom099/factory/inventory/internal/service/model"
)

func (r *repository) Create(ctx context.Context, dto model.Part) (string, error) {
	uuid := uuid.New().String()
	dto.Uuid = uuid
	dto.CreatedAt = time.Now()

	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[uuid] = converter.ToRepoPart(dto)

	return uuid, nil
}
