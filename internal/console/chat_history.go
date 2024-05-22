package console

import (
	"fmt"

	"github.com/Tel3scop/chat-client/internal/model"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// DefaultHistoryMessages количество сообщений в истории по умолчанию
const DefaultHistoryMessages = 20

func (c *Console) runHistory(_ *cobra.Command, _ []string) {
	fmt.Println(otherMessage("Команда истории чата.\n"))
	if c.currentChat.ID == 0 {
		fmt.Println(color.RedString("Вы должны подключиться к чату, чтобы использовать эту команду"))
		return
	}

	c.whereIAm()
	messages, err := c.getChatHistory(DefaultHistoryMessages)
	if err != nil {
		fmt.Println(color.RedString("Ошибка при получении истории чата: %s", err))
		return
	}

	c.displayMessages(messages)
}

func (c *Console) getChatHistory(count int64) ([]model.ChatMessage, error) {
	messages, err := c.chatService.GetMessages(c.ctx, c.currentChat.ID, count)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (c *Console) displayMessages(messages []model.ChatMessage) {
	if len(messages) == 0 {
		fmt.Println("Нет сообщений в истории чата.")
		return
	}

	fmt.Println("История сообщений:")
	for i, msg := range messages {
		fmt.Printf("[%d] [%v] %s: %s\n", i+1, msg.Timestamp.AsTime().Format("02.01.2006 15:04"), msg.From, msg.Text)
	}
}
