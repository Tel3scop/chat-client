package auth

import (
	"context"
)

func (s *serv) Login(ctx context.Context, username, password string) error {
	err := s.authClient.Login(ctx, username, password)
	if err != nil {
		return err
	}
	return nil
}
