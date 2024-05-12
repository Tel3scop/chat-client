package console

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

var myMessage = color.New(color.Italic, color.FgGreen).SprintFunc()
var otherMessage = color.New(color.Bold, color.FgBlue).SprintFunc()
var menuTitle = color.New(color.Bold, color.FgGreen).SprintFunc()

// Запрашивает ввод пользователя
func (c *Console) promptInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	return input, err
}

// Обрабатывает ошибки ввода/вывода
func (c *Console) handleIOError(err error, message string) {
	fmt.Println(color.RedString(message+": %s", err))
}
