package console

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var chatCmd *cobra.Command

// Общий цикл для управления чатами
func (c *Console) runChatMenu(_ *cobra.Command, _ []string) {
	if err := c.checkAuth(); err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
		return
	}

	fmt.Println(otherMessage("Режим чата.\n"))
	fmt.Printf("Доступные команды: %s\n", helpChatText)

	for {
		fmt.Println(color.GreenString("Вы находитесь в меню управления чатами."))
		input, err := c.promptInput("Введите команду чата: ")
		if err != nil {
			fmt.Printf("%s\n", color.RedString("Ошибка: %s", err))
			continue
		}

		input = strings.TrimSpace(input)
		c.cmd.SetArgs(strings.Split(input, " "))
		if err = c.cmd.Execute(); err != nil {
			fmt.Println(err)
			if input == Exit {
				break
			}
		}
	}
}

// Обработка команд чата
func (c *Console) executeChatCommand(input string) error {
	switch input {
	case Exit:
		return errors.New("exit")
	case Show:
		c.showMyChats()
	case Connect:
		c.runConnectChat(c.cmd, nil)
	case Create:
		c.runCreateChat(c.cmd, nil)
	case History:
		c.runHistory(c.cmd, nil)
	default:
		return fmt.Errorf(color.RedString("Команда не найдена: %s", input))
	}
	return nil
}

// Отображение доступных чатов
func (c *Console) showMyChats() {
	chats, err := c.chatService.GetChats(c.ctx, c.username)
	if err != nil {
		fmt.Println(color.RedString("Ошибка при получении списка доступных чатов: %s", err))
		return
	}
	if len(chats) > 0 {
		fmt.Println("Доступные чаты:")
		header := fmt.Sprintf("%s | %s | %s", color.YellowString("ID"), color.CyanString("Название"), color.GreenString("Участники"))
		fmt.Println(header)
		fmt.Println(strings.Repeat("-", len(header)))
		for _, chat := range chats {
			myChats[chat.ID] = chat

			line := fmt.Sprintf("%s | %s | %s",
				color.YellowString(fmt.Sprint(chat.ID)),
				color.CyanString(chat.Name),
				color.GreenString(strings.Join(chat.Members, ", ")))
			fmt.Println(line)
		}
	} else {
		fmt.Println(color.RedString("Нет доступных чатов."))
	}
}
