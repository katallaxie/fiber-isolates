// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package isolates

import (
	"github.com/gofiber/fiber/v2"
	"github.com/katallaxie/v8go-polyfills/listener"

	v8 "github.com/katallaxie/v8go"
)

// Injector ...
type Injector func(*fiber.Ctx, *v8.Isolate, *v8.ObjectTemplate) error

// Config is the config for the isolates middleware
type Config struct {
	// Filter defines a function to skip the middleware.
	// Optional. Default: nil
	Filter func(*fiber.Ctx) bool

	// Injetion defines a function to inject the context into the isolate.
	// Optional. Default: nil
	Injetion []Injector

	// Next defines a function to skip this middleware
	Next func(*fiber.Ctx) bool
}

// New ...
func New(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set default config
		_ = configDefault(config)

		iso := v8.NewIsolate()
		global := v8.NewObjectTemplate(iso)

		defer iso.Dispose()

		for _, inject := range config.Injetion {
			if err := inject(c, iso, global); err != nil {
				return c.SendStatus(fiber.StatusBadRequest)
			}
		}

		in := make(chan *v8.Object)
		out := make(chan *v8.Value)

		l := listener.New(listener.WithEvents("request", in, out))

		if err := l.Inject(iso, global); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		ctx := v8.NewContext(iso, global)
		defer ctx.Close()

		_, err := ctx.RunScript("addListener('request', event => { return event.sourceIP === '127.0.0.1' })", "listener.js")
		if err != nil {

			return c.SendStatus(fiber.StatusBadRequest)
		}

		obj, err := newContextObject(ctx)
		if err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}

		in <- obj

		v := <-out
		defer func() { v.Release() }()

		if v.IsBoolean() {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.SendStatus(fiber.StatusUnauthorized)
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
