package chat

import (
	"context"

	"github.com/Tel3scop/chat-client/internal/model"
)

func (s *serv) GetMessages(ctx context.Context, chatID, count int64) ([]model.ChatMessage, error) {
	messages, err := s.chatClient.GetMessages(ctx, chatID, count)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
