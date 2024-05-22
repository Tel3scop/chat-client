package chat

import (
	"context"
	"time"

	"github.com/Tel3scop/chat-client/internal/model"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *serv) SendMessage(ctx context.Context, ID int64, from string, text string, timestamp *time.Time) error {
	request := model.ChatMessage{
		ChatID: ID,
		From:   from,
		Text:   text,
	}
	request.Timestamp = timestamppb.Now()
	if timestamp != nil {
		request.Timestamp = timestamppb.New(*timestamp)
	}

	err := s.chatClient.SendMessage(ctx, request)
	if err != nil {
		return err
	}

	return nil
}
