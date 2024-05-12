package app

import (
	"context"
	"log"

	"github.com/Tel3scop/chat-client/internal/config"
	"github.com/Tel3scop/chat-client/internal/connector/auth"
	"github.com/Tel3scop/chat-client/internal/connector/chat"
	"github.com/Tel3scop/chat-client/internal/console"
	"github.com/Tel3scop/chat-client/internal/cron"
	"github.com/Tel3scop/chat-client/internal/service"
	authService "github.com/Tel3scop/chat-client/internal/service/auth"
	chatService "github.com/Tel3scop/chat-client/internal/service/chat"
)

type serviceProvider struct {
	authService service.AuthService
	chatService service.ChatService
	config      *config.Config
	authClient  *auth.Client
	chatClient  *chat.Client
	console     *console.Console
	cron        *cron.Cron
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("failed to get config: %s", err.Error())
		}
		s.config = cfg
	}
	return s.config
}

func (s *serviceProvider) ChatClient() *chat.Client {

	if s.chatClient == nil {
		var err error
		s.chatClient, err = chat.New(s.Config().ChatService.Host, s.Config().ChatService.Port)
		if err != nil {
			log.Fatal("Failed to create auth client", err.Error())
		}
	}
	return s.chatClient
}

func (s *serviceProvider) AuthClient() *auth.Client {

	if s.authClient == nil {
		var err error
		s.authClient, err = auth.New(s.Config().AuthService.Host, s.Config().AuthService.Port)
		if err != nil {
			log.Fatal("Failed to create auth client", err.Error())
		}
	}
	return s.authClient
}

func (s *serviceProvider) ChatService() service.ChatService {
	if s.chatService == nil {
		s.chatService = chatService.NewService(
			s.ChatClient(),
		)
	}

	return s.chatService
}

func (s *serviceProvider) AuthService() service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(
			s.AuthClient(),
		)
	}

	return s.authService
}

func (s *serviceProvider) Cron() *cron.Cron {
	if s.cron == nil {
		s.cron = cron.NewCron(
			s.AuthClient(),
			s.Config(),
		)
	}

	return s.cron
}

func (s *serviceProvider) Console(ctx context.Context) *console.Console {
	if s.console == nil {
		s.console = console.NewConsole(ctx,
			s.AuthService(),
			s.ChatService(),
		)
	}

	return s.console
}
