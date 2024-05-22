package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Tel3scop/chat-client/internal/app"
	"github.com/Tel3scop/chat-client/internal/cron"
	"github.com/fatih/color"
)

func main() {
	flag.Parse()
	fmt.Println(color.GreenString("Добро пожаловать в утилиту чат-клиента!"))
	ctx := context.Background()
	newApp, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("Не удалось запустить приложение: %s", err)
	}

	defer cron.StopCron()
	if newApp.Run(ctx) != nil {
		os.Exit(0)
	}

	os.Exit(1)

}
