package session

import "context"

func (r *repository) AddSessionToUserSet(ctx context.Context, userID, sessionID string) error {
	// добавляет новую сессию пользователя в множество его сессий
	// чтобы можно было управлять сессиями
	key := r.getUserSessionsKey(userID)

	return r.cache.SAdd(ctx, key, sessionID)
}
