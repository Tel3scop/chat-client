package chat

import (
	"context"
)

func (s *serv) Create(ctx context.Context, name string, members []string) (int64, error) {
	chatID, err := s.chatClient.Create(ctx, name, members)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
