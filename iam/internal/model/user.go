package model

import "time"

type User struct {
	ID        string
	Login     string
	Email     string
	Hash      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}
