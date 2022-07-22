package middleware

import (
	"github.com/carlosrodriguesf/bank-api/pkg/api/middleware/auth"
	"github.com/carlosrodriguesf/bank-api/pkg/app"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
)

type (
	Options struct {
		Logger logger.Logger
		Apps   app.Container
	}
	Container interface {
		Auth() auth.Middleware
	}
	container struct {
		auth auth.Middleware
	}
)

func NewContainer(opts Options) Container {
	return &container{
		auth: auth.NewMiddleware(auth.Options{
			Logger: opts.Logger,
			Apps:   opts.Apps,
		}),
	}
}

func (c *container) Auth() auth.Middleware {
	return c.auth
}
