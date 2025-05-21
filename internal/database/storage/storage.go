package storage

import "context"

type Storage interface {
	AddChatID(ctx context.Context, chatID string) error
	RemoveChatID(ctx context.Context, chatID string) error
	GetChatIDs(ctx context.Context) ([]string, error)
	Close() error
}
