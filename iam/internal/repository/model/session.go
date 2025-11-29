package model

// todo - решить, какие поля храним в редисе

type SessionRedisView struct {
	UUID         string  `redis:"uuid"`
	ObservedAtNs *int64  `redis:"observed_at,omitempty"`
	Description  string  `redis:"description"`
	Color        *string `redis:"color,omitempty"`
	Sound        *bool   `redis:"sound,omitempty"`
	Duration     *int32  `redis:"duration_seconds,omitempty"`
	CreatedAtNs  int64   `redis:"created_at"`
	UpdatedAtNs  *int64  `redis:"updated_at,omitempty"`
	DeletedAtNs  *int64  `redis:"deleted_at,omitempty"`
}
