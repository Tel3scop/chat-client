package console

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func (c *Console) runConnectChat(_ *cobra.Command, _ []string) {
	if err := c.checkAuth(); err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
		return
	}

	chatID, err := c.promptForChatID()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := c.connectToChat(chatID); err != nil {
		fmt.Println(err)
		return
	}

	c.handleChatSession(chatID)
}

func (c *Console) promptForChatID() (int64, error) {
	c.whereIAm()
	input, err := c.promptInput("Введите ID чата: ")
	if err != nil {
		return 0, fmt.Errorf("ошибка при вводе ID чата: %v", err)
	}

	chatID, err := strconv.ParseInt(strings.TrimSpace(input), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("введен некорректный ID чата. Необходимо ввести число")
	}

	return chatID, nil
}

func (c *Console) connectToChat(chatID int64) error {
	foundedChat, ok := myChats[chatID]
	if !ok {
		return errors.New(fmt.Sprintf("не удалось найти чат %d", chatID))
	}

	c.currentChat = foundedChat
	fmt.Printf("Connected to chat: %s\n", foundedChat.Name)

	return nil
}

func (c *Console) handleChatSession(chatID int64) {
	c.connectStream(c.ctx, chatID)

	for {
		message, err := c.promptInput("Введите текст сообщения: ")
		if err != nil {
			fmt.Println("Ошибка при вводе сообщения:", err)
			return
		}

		if strings.TrimSpace(message) == Exit {
			fmt.Println("Выход из чата.")
			break
		}

		if err := c.sendMessage(c.ctx, chatID, message); err != nil {
			fmt.Println(err)
			continue
		}
	}
}

func (c *Console) sendMessage(ctx context.Context, chatID int64, message string) error {
	timestamp := time.Now()
	if err := c.chatService.SendMessage(ctx, chatID, c.username, message, &timestamp); err != nil {
		return fmt.Errorf("ошибка при отправке сообщения в чат: %v", err)
	}
	fmt.Printf("[%v][%s]: %s\n", color.YellowString(timestamp.Format("02.01.2006 15:04")), color.BlueString(c.username), myMessage(message))
	return nil
}

func (c *Console) connectStream(ctx context.Context, chatID int64) {
	fmt.Println(otherMessage("Команда подключения к чату.\n"))
	stream, err := c.chatService.ConnectChat(ctx, chatID, c.username)
	if err != nil {
		logger.Error("can't connect chatElement", zap.Error(err))
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				fmt.Println("failed to receive message from stream: ", errRecv)
				logger.Error("failed to receive message from stream: ", zap.Error(errRecv))
				return
			}

			fmt.Printf("[%v][%s]: %s\n",
				color.YellowString(message.GetCreatedAt().AsTime().Format("02.01.2006 15:04")),
				color.BlueString(message.GetFrom()),
				otherMessage(message.GetText()),
			)
		}
	}()
}
