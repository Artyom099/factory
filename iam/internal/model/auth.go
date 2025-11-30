package model

import "time"

type Session struct {
	ID        string
	UserID    string
	CreatedAt time.Time
	UpdatedAt *time.Time
	ExpiredAt time.Time
}
