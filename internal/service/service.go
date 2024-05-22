package service

import (
	"context"
	"time"

	"github.com/Tel3scop/chat-client/internal/model"
	"github.com/Tel3scop/chat-client/internal/pkg/chat_v1"
)

// AuthService интерфейс для использования в сервисе
type AuthService interface {
	Login(ctx context.Context, username, password string) error
}

// ChatService интерфейс для использования в сервисе
type ChatService interface {
	Create(ctx context.Context, name string, members []string) (int64, error)
	SendMessage(ctx context.Context, ID int64, from string, text string, timestamp *time.Time) error
	ConnectChat(ctx context.Context, chatID int64, username string) (chat_v1.ChatV1_ConnectChatClient, error)
	GetChats(ctx context.Context, name string) ([]model.Chat, error)
	GetMessages(ctx context.Context, chatID, count int64) ([]model.ChatMessage, error)
}
