package model

type SessionRedisView struct {
	UUID        string `redis:"uuid"`
	UserID      string `redis:"user_id"`
	CreatedAtNs int64  `redis:"created_at"`
	UpdatedAtNs *int64 `redis:"updated_at,omitempty"`
	ExpiredAtNs int64  `redis:"expired_at"`
}
