package auth

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Tel3scop/auth/pkg/auth_v1"
	"github.com/Tel3scop/chat-client/internal/connector"
	"github.com/Tel3scop/helpers/logger"
	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Client экземпляр
type Client struct {
	client auth_v1.AuthV1Client
}

var _ connector.AuthClient = (*Client)(nil)

// TokenStruct структура для управления токеном
type TokenStruct struct {
	value string
	mu    sync.Mutex
}

var refreshToken TokenStruct
var accessToken TokenStruct

// New создает новый экземпляр клиента
func New(host string, port int64) (*Client, error) {
	refreshToken = TokenStruct{}
	accessToken = TokenStruct{}

	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial GRPC client: %v", err)
	}

	return &Client{
		client: auth_v1.NewAuthV1Client(conn),
	}, nil
}

// Login метод авторизации в сервисе
func (c *Client) Login(ctx context.Context, username, password string) error {
	response, err := c.client.Login(ctx, &auth_v1.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return err
	}
	refreshToken.mu.Lock()
	defer refreshToken.mu.Unlock()

	refreshToken.value = response.RefreshToken

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err = c.SetAccessToken(ctx)
		if err != nil {
			logger.Error("failed to set access token", zap.Error(err))
		}
		logger.Info("access token set successfully")
	}()
	return nil
}

// SetRefreshToken метод обновления refresh-токена
func (c *Client) SetRefreshToken(ctx context.Context) error {
	token := GetRefreshToken()
	response, err := c.client.GetRefreshToken(ctx, &auth_v1.GetRefreshTokenRequest{OldRefreshToken: token})
	if err != nil {
		return err
	}
	logger.Info("updated refresh token")
	refreshToken.mu.Lock()
	defer refreshToken.mu.Unlock()

	refreshToken.value = response.RefreshToken
	return nil
}

// SetAccessToken метод получения access-токена
func (c *Client) SetAccessToken(ctx context.Context) error {
	token := GetRefreshToken()
	response, err := c.client.GetAccessToken(ctx, &auth_v1.GetAccessTokenRequest{RefreshToken: token})
	if err != nil {
		return err
	}
	logger.Info("updated access token")
	accessToken.mu.Lock()
	defer accessToken.mu.Unlock()

	accessToken.value = response.AccessToken
	return nil
}

// GetRefreshToken Вернуть refresh-токен
func GetRefreshToken() string {
	refreshToken.mu.Lock()
	defer refreshToken.mu.Unlock()
	return refreshToken.value
}

// GetAccessToken Вернуть access-токен
func GetAccessToken() string {
	accessToken.mu.Lock()
	defer accessToken.mu.Unlock()
	return accessToken.value
}
