package swagger

import (
	"github.com/carlosrodriguesf/bank-api/docs"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// Options ...
type Options struct {
	Logger  logger.Logger
	Echo    *echo.Echo
	Title   string
	Version string
}

// Register group item check
func Register(opts Options) {
	log := opts.Logger.WithPreffix("api.swagger")

	docs.SwaggerInfo.Title = opts.Title
	docs.SwaggerInfo.Version = opts.Version
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	opts.Echo.GET("/docs", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/docs/index.html")
	})

	opts.Echo.GET("/docs/*", func(c echo.Context) error {
		return echoSwagger.WrapHandler(c)
	})

	log.Info("started")
}
