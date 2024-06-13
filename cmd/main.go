package main

import (
	"log/slog"
	"os"
	"os/signal"

	"github.com/Richtermnd/RichterAuth/internal/application"
	"github.com/Richtermnd/RichterAuth/internal/config"
)

func main() {
	config.LoadConfig()
	log := slog.Default()
	app := application.New(log)
	go app.Run()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	<-ch
	app.Stop()
}
