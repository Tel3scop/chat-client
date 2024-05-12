package console

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helpText = `Справка по командам главного меню: 
/help - вызов справки
/chat - переход в раздел управления чатами. Чтобы начать общение - перейдите в данный раздел и выберите нужный чат
/exit - выйти из раздела помощи / завершение программы
------------------------
`
var helpChatText = `
Справка по командам меню чата:
/create - меню создания нового чата
/show - показать доступные чаты
/connect - подключение к соответствующему чату
/history (доступна, если Вы находитесь в чате) - показать последние 20 сообщений чата
/exit - выйди из чата
`

var fullHelp = helpText + helpChatText

func (c *Console) runHelp(cmd *cobra.Command, args []string) {
	fmt.Printf(fullHelp)
	return
}
