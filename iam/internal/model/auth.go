package model

import "time"

type Session struct {
	ID        string
	UserID    string
	Login     string
	CreatedAt time.Time
	UpdatedAt *time.Time
	ExpiredAt time.Time
}
