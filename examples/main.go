package main

import (
	isolates "github.com/katallaxie/fiber-isolates"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{})

	// Add logger
	app.Use(logger.New())

	app.Get("/index", isolates.New(isolates.Config{}))

	if err := app.Listen(":8080"); err != nil {
		panic(err)
	}
}
