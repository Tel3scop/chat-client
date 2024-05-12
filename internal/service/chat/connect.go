package chat

import (
	"context"

	"github.com/Tel3scop/chat-client/internal/pkg/chat_v1"
)

func (s *serv) ConnectChat(ctx context.Context, chatID int64, username string) (chat_v1.ChatV1_ConnectChatClient, error) {
	stream, err := s.chatClient.ConnectChat(ctx, chatID, username)
	if err != nil {
		return nil, err
	}

	return stream, nil
}
