package chat

import (
	"context"

	"github.com/Tel3scop/chat-client/internal/connector"
	"github.com/Tel3scop/chat-client/internal/connector/auth"
	"github.com/Tel3scop/chat-client/internal/model"
	"github.com/Tel3scop/chat-client/internal/pkg/chat_v1"

	//todo: поменять на "github.com/Tel3scop/chat-server/pkg/chat_v1" когда смердим ветку week8
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Client экземпляр
type Client struct {
	client chat_v1.ChatV1Client
}

var _ connector.ChatClient = (*Client)(nil)

// New создает новый экземпляр клиента
func New(conn *grpc.ClientConn) (*Client, error) {
	return &Client{
		client: chat_v1.NewChatV1Client(conn),
	}, nil
}

// Create метод создания нового чата
func (c *Client) Create(ctx context.Context, name string, members []string) (int64, error) {

	response, err := c.client.Create(setAccessContext(ctx), &chat_v1.CreateRequest{
		Usernames: members,
		Name:      name,
	})
	if err != nil {
		return 0, err
	}

	return response.GetId(), nil
}

// SendMessage метод отправки сообщения
func (c *Client) SendMessage(ctx context.Context, request model.ChatMessage) error {
	_, err := c.client.SendMessage(setAccessContext(ctx), &chat_v1.SendMessageRequest{
		ChatId:    request.ChatID,
		From:      request.From,
		Text:      request.Text,
		Timestamp: request.Timestamp,
	})
	if err != nil {
		return err
	}

	return nil
}

// ConnectChat метод подключения к чату
func (c *Client) ConnectChat(ctx context.Context, chatID int64, username string) (chat_v1.ChatV1_ConnectChatClient, error) {
	stream, err := c.client.ConnectChat(setAccessContext(ctx), &chat_v1.ConnectChatRequest{
		ChatId:   chatID,
		Username: username,
	})
	if err != nil {
		return nil, err
	}

	return stream, nil
}

// GetChats метод получения чатов
func (c *Client) GetChats(ctx context.Context, username string) ([]model.Chat, error) {
	chats, err := c.client.GetChats(setAccessContext(ctx), &chat_v1.GetChatsRequest{Name: username})
	if err != nil {
		return nil, err
	}
	result := make([]model.Chat, 0, len(chats.Chats))
	for _, chat := range chats.Chats {
		result = append(result, model.Chat{
			ID:      chat.Id,
			Name:    chat.Name,
			Members: chat.Usernames,
		})
	}

	return result, nil
}

// GetMessages метод получения сообщений
func (c *Client) GetMessages(ctx context.Context, chatID, count int64) ([]model.ChatMessage, error) {
	messages, err := c.client.GetMessages(setAccessContext(ctx), &chat_v1.GetMessagesRequest{
		ChatId: chatID,
		Count:  count,
	})
	if err != nil {
		return nil, err
	}
	result := make([]model.ChatMessage, 0, len(messages.Messages))
	for _, msg := range messages.Messages {
		result = append(result, model.ChatMessage{
			ChatID:    chatID,
			From:      msg.From,
			Text:      msg.Text,
			Timestamp: msg.CreatedAt,
		})
	}

	return result, nil
}

func setAccessContext(ctx context.Context) context.Context {
	md := metadata.New(map[string]string{"Authorization": "Bearer " + auth.GetAccessToken()})
	return metadata.NewOutgoingContext(ctx, md)
}
