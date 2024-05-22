package console

import (
	"context"
	"fmt"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"golang.org/x/term"
)

func (c *Console) attemptLogin(ctx context.Context) bool {
	for i := 0; i < MaxLoginTries; i++ {
		usernameValue, err := c.promptInput("Введите имя пользователя: ")
		if err != nil {
			c.handleIOError(err, "Ошибка при вводе имени пользователя")
			return false
		}
		c.username = strings.TrimSpace(usernameValue)

		password, err := c.promptPassword("Введите пароль: ")
		if err != nil {
			c.handleIOError(err, "Ошибка при вводе пароля")
			return false
		}

		if err = c.authService.Login(ctx, c.username, password); err == nil {
			time.Sleep(time.Second)
			return true
		}

		fmt.Println(color.RedString("Неверное имя пользователя или пароль"))
	}
	return false
}

// Запрашивает пароль без отображения ввода
func (c *Console) promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(syscall.Stdin)
	if err != nil {
		return "", err
	}

	return string(bytePassword), nil
}

// Проверка авторизации пользователя
func (c *Console) checkAuth() error {
	if !c.isLoggedIn {
		return fmt.Errorf("not authenticated")
	}
	return nil
}
