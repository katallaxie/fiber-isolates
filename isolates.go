// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package isolates

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ionos-cloud/v8go-polyfills/console"
	"github.com/ionos-cloud/v8go-polyfills/listener"

	v8 "rogchap.com/v8go"
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
		iso := v8.NewIsolate()
		global := v8.NewObjectTemplate(iso)

		events := make(chan *v8.Object)

		if err := listener.AddTo(iso, global, listener.WithEvents("request", events)); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		ctx := v8.NewContext(iso, global)

		_, err := ctx.RunScript("addListener('request', event => { return event.sourceIP === '127.0.0.1' })", "listener.js")
		if err != nil {

			return c.SendStatus(fiber.StatusBadRequest)
		}

		if err := console.AddTo(ctx); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		obj, err := newContextObject(ctx)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		events <- obj

		return c.SendStatus(fiber.StatusOK)
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

func newContextObject(ctx *v8.Context) (*v8.Object, error) {
	iso := ctx.Isolate()
	obj := v8.NewObjectTemplate(iso)

	resObj, err := obj.NewInstance(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range []struct {
		Key string
		Val interface{}
	}{
		{Key: "sourceIP", Val: "127.0.0.1"},
	} {
		if err := resObj.Set(v.Key, v.Val); err != nil {
			return nil, err
		}
	}

	return resObj, nil
}
