package http

import "context"

type ITelegramClient interface {
	SendMessage(ctx context.Context, chatID int64, text string) error
}
