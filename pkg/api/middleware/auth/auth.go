package auth

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	"github.com/carlosrodriguesf/bank-api/pkg/app"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/labstack/echo/v4"
)

type (
	Options struct {
		Logger logger.Logger
		Apps   app.Container
	}
	Middleware interface {
		Private(next echo.HandlerFunc) echo.HandlerFunc
	}
	middlewareImpl struct {
		logger logger.Logger
		apps   app.Container
	}
)

func NewMiddleware(opts Options) Middleware {
	return &middlewareImpl{
		logger: opts.Logger.WithLocation().WithPreffix("api.middleware.auth"),
		apps:   opts.Apps,
	}
}

func (a *middlewareImpl) Private(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := getTokenFromRequest(c.Request())
		if token == "" {
			return apierror.ErrUnauthorized
		}

		ctx := c.Request().Context()
		session, err := a.apps.Auth().GetSessionByToken(ctx, token)
		if err != nil {
			if err := apierror.Get(err, errorMap); err != nil {
				return err
			}
			a.logger.Error(err)
			return apierror.ErrInternal
		}

		ctx = model.SetSessionOnContext(ctx, session)
		c.SetRequest(c.Request().WithContext(ctx))
		return next(c)
	}
}
