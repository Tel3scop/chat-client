package chat

import (
	"context"

	"github.com/Tel3scop/chat-client/internal/model"
)

func (s *serv) GetChats(ctx context.Context, name string) ([]model.Chat, error) {
	chats, err := s.chatClient.GetChats(ctx, name)
	if err != nil {
		return nil, err
	}

	return chats, nil
}
