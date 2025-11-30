package model

import "time"

type RepoNotificationMethod struct {
	ID           string `db:"id"`
	UserID       string `db:"user_id"`
	ProviderName string `db:"provider_name"`
	Target       string `db:"target"`
}

type RepoUser struct {
	ID                  string                   `db:"id"`
	Login               string                   `db:"login"`
	Email               string                   `db:"email"`
	Hash                string                   `db:"hash"`
	NotificationMethods []RepoNotificationMethod `db:"-"`
	CreatedAt           time.Time                `db:"created_at"`
	UpdatedAt           *time.Time               `db:"updated_at"`
}
