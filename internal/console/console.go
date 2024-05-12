package console

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Tel3scop/chat-client/internal/model"
	"github.com/Tel3scop/chat-client/internal/service"
	"github.com/Tel3scop/helpers/logger"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Exit выход
const Exit = "exit"

// Help раздел помощи
const Help = "help"

// Show показать доступные чаты
const Show = "show"

// Chat раздел чатов
const Chat = "chat"

// Create создание нового чата
const Create = "create"

// Connect подключение к чату
const Connect = "connect"

// History показать историю чата
const History = "history"

const MaxLoginTries = 3

// Console структура для работы с консольными командами
type Console struct {
	ctx         context.Context
	authService service.AuthService
	chatService service.ChatService
	currentChat model.Chat
	isLoggedIn  bool
	username    string
	cmd         *cobra.Command
	reader      *bufio.Reader
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "chat-client",
	Short: "Чат-клиент",
	Long:  "Чат-клиент позволяет через консоль создавать чаты с друзьями",
}

// myChats список чатов пользователя
var myChats = make(map[int64]model.Chat)

func (c *Console) whereIAm() {
	fmt.Printf(color.GreenString("\n\n Уважаемый %s! ", c.username))
	if c.currentChat.ID != 0 {
		fmt.Printf(color.GreenString("Вы находитесь в чате [%d] %s", c.currentChat.ID, c.currentChat.Name))
	} else {
		fmt.Printf(color.YellowString("Вы не покдлючены к чату"))
	}
	fmt.Println()
}

// NewConsole новый экземпляр структуры Console
func NewConsole(ctx context.Context, authService service.AuthService, chatService service.ChatService) *Console {
	console := &Console{
		ctx:         ctx,
		authService: authService,
		chatService: chatService,
		reader:      bufio.NewReader(os.Stdin),
	}

	console.setupCommands()

	return console
}

func (c *Console) setupCommands() {
	c.cmd = &cobra.Command{
		Use:   "chat-client",
		Short: "Чат-клиент",
		Run:   c.runRoot,
	}

	c.cmd.AddCommand(
		&cobra.Command{
			Use:   Chat,
			Short: "Управление чатами",
			Run:   c.runChatMenu,
		},
		&cobra.Command{
			Use:   Connect,
			Short: "Подключиться к чату",
			Run:   c.runConnectChat,
		},
		&cobra.Command{
			Use:   Create,
			Short: "Создать новый чат",
			Run:   c.runCreateChat,
		},
		&cobra.Command{
			Use:   History,
			Short: "Посмотреть историю по чату",
			Run:   c.runHistory,
		},
		&cobra.Command{
			Use:   Show,
			Short: "Посмотреть историю по чату",
			Run:   c.runShow,
		},
		&cobra.Command{
			Use:   Help,
			Short: "Помощь",
			Run:   c.runHelp,
		},
	)
}

func (c *Console) runRoot(cmd *cobra.Command, args []string) {
	if err := c.checkAuth(); err != nil {
		logger.Error("Authentication required", zap.Error(err))
		return
	}

	for {
		fmt.Print("Enter command: ")
		input, _ := c.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		switch input {
		case Exit:
			fmt.Println("Exiting...")
			return
		default:
			c.cmd.SetArgs(strings.Split(input, " "))
			c.cmd.Execute()
		}
	}
}

func (c *Console) Run() {
	ctx := context.Background()

	if !c.attemptLogin(ctx) {
		fmt.Println(color.RedString("Превышено максимальное количество попыток! Выход"))
		os.Exit(1)
	}

	fmt.Println("Вы успешно авторизовались")
	c.isLoggedIn = true

	c.showMyChats()
	c.executeCommandLoop()
}

// Цикл выполнения команд
func (c *Console) executeCommandLoop() {
	for {
		c.whereIAm()
		fmt.Println(color.GreenString("\nВы находитесь в главном меню."))
		input, err := c.promptInput("Введите команду: ")
		if err != nil {
			c.handleIOError(err, "Ошибка при чтении команды")
			continue
		}

		input = strings.TrimSpace(input)
		if input == Exit {
			fmt.Println(color.GreenString("Возвращайтесь скорее!"))
			os.Exit(0)
		}

		c.cmd.SetArgs(strings.Split(input, " "))
		c.cmd.Execute()
	}
}

func (c *Console) execute(args []string) error {
	rootCmd.SetArgs(args)
	err := rootCmd.Execute()
	if err != nil {
		return err
	}

	return nil
}
