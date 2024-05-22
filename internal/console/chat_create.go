package console

import (
	"context"
	"fmt"

	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func (c *Console) runCreateChat(_ *cobra.Command, _ []string) {
	c.whereIAm()

	fmt.Println(otherMessage("Команда создания чата.\n"))
	name, err := c.promptInput("Введите название чата: ")
	if err != nil {
		fmt.Println("Ошибка при вводе названия чата:", err)
		return
	}

	members, err := c.promptChatMembers()
	if err != nil {
		fmt.Println("Ошибка при вводе имен участников:", err)
		return
	}

	chatID, err := c.createChat(c.ctx, name, members)
	if err != nil {
		fmt.Println(color.RedString("Ошибка при создании чата: %s", err))
		return
	}

	fmt.Printf("Чат успешно создан:\nID: %d\nНазвание: %s\nСостав: %v\n", chatID, name, members)
}

func (c *Console) promptChatMembers() ([]string, error) {
	fmt.Println("Введите имена участников чата (через запятую): ")
	input, err := c.promptInput("")
	if err != nil {
		return nil, err
	}

	members := strings.Split(strings.TrimSpace(input), ",")
	// Добавляем текущего пользователя автоматически
	members = append(members, c.username)

	return members, nil
}

func (c *Console) createChat(ctx context.Context, name string, members []string) (int64, error) {
	chatID, err := c.chatService.Create(ctx, name, members)
	if err != nil {
		return 0, err
	}

	return chatID, nil
}
