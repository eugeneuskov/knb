package main

import (
	"knb/app"
	"knb/app/config"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	envFilePath = ".env"
)

func main() {
	appConfig, err := new(config.Config).Init(envFilePath)
	if err != nil {
		log.Fatalf("Failed initializing config: %s\n", err.Error())
		return
	}

	println("App starting...")
	application := app.NewApplication(appConfig)
	application.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	println("\nApp shutting down...")

	application.Shutdown()
}
