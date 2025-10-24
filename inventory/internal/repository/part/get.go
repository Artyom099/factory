package part

import (
	"context"

	"github.com/Artyom099/factory/inventory/internal/repository/converter"
	servModel "github.com/Artyom099/factory/inventory/internal/service/model"
)

func (r *repository) Get(ctx context.Context, uuid string) (servModel.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[uuid]
	if !ok {
		return servModel.Part{}, servModel.ErrPartNotFound
	}

	return converter.ToModelPart(part), nil
}
