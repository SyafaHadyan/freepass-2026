package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SyafaHadyan/freepass-2026/internal/bootstrap"
)

func main() {
	app := bootstrap.Start()

	go func() {
		err := app.App.Fiber.Listen(fmt.Sprintf(":%d", app.Config.AppPort))
		if err != nil {
			log.Panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Println("gracefully shutting down")

	err := app.App.Fiber.Shutdown()
	if err != nil {
		log.Println("graceful shutdown failed")
	}

	log.Println("shutdown complete")
}
