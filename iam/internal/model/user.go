package model

import "time"

type User struct {
	ID                  string
	Login               string
	Email               string
	Password            string
	Hash                string
	NotificationMethods []NotificationMethod // Каналы уведомлений - telegram, email, push и т.д.
	CreatedAt           time.Time
	UpdatedAt           *time.Time
}

type NotificationMethod struct {
	ProviderName string
	Target       string
}
