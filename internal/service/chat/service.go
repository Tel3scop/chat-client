package chat

import (
	"github.com/Tel3scop/chat-client/internal/connector/chat"
	"github.com/Tel3scop/chat-client/internal/service"
)

type serv struct {
	chatClient *chat.Client
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	chatClient *chat.Client,
) service.ChatService {
	return &serv{
		chatClient: chatClient,
	}
}
