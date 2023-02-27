// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package isolates

import (
	"github.com/gofiber/fiber/v2"
)

// Config is the config for the isolates middleware
type Config struct {
	// Filter defines a function to skip the middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// Next defines a function to skip this middleware
	Next func(*fiber.Ctx) bool
}

// New ...
func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNoContent)
	}
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	return cfg
}
