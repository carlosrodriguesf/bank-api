package api

import (
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	v1 "github.com/carlosrodriguesf/bank-api/pkg/api/v1"
	"github.com/labstack/echo/v4"
)

// Register configure swagger
// Api use bearer token in private requests
// @securityDefinitions.apikey UserToken
// @in header
// @name Authorization
func Register(e *echo.Echo, opts apimodel.Options) {
	v1.Register(e.Group("/api"), opts)

	opts.Logger.WithPreffix("api").Info("registered")
}
