package auth

import (
	"github.com/Tel3scop/chat-client/internal/connector/auth"
	"github.com/Tel3scop/chat-client/internal/service"
)

type serv struct {
	authClient *auth.Client
}

// NewService функция возвращает новый сервис пользователя
func NewService(
	authClient *auth.Client,
) service.AuthService {
	return &serv{
		authClient: authClient,
	}
}
