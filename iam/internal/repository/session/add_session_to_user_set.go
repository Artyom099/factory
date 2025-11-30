package session

import "context"

func (r *repository) AddSessionToUserSet(ctx context.Context, userID, sessionID string) error {
	key := r.getUserSessionsKey(userID)

	return r.cache.SAdd(ctx, key, sessionID)
}
