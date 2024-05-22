package connector

import (
	"context"

	"github.com/Tel3scop/chat-client/internal/model"
)

// AuthClient is a client for authentication service.
type AuthClient interface {
	Login(ctx context.Context, username, password string) error
	SetRefreshToken(ctx context.Context) error
	SetAccessToken(ctx context.Context) error
}

// ChatClient is a client for chat service.
type ChatClient interface {
	Create(ctx context.Context, name string, members []string) (int64, error)
	SendMessage(ctx context.Context, request model.ChatMessage) error
	GetChats(ctx context.Context, username string) ([]model.Chat, error)
}
