package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()

    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "status": "ok",
            "service": "matchaciee-api",
        })
    })

    app.Get("/", func (c *fiber.Ctx) error {
        return c.SendString("Hello, World!")
    })

    // Channel to listen for interrupt signals
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)

    // Start server in a goroutine
    go func() {
        if err := app.Listen(":8080"); err != nil {
            log.Panic(err)
        }
    }()

    // Block until we receive a signal
    <-c
    log.Println("Gracefully shutting down...")
    _ = app.Shutdown()
    log.Println("Server stopped")
}