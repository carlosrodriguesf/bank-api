package v1

import (
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group, opts apimodel.Options) {
	opts.Logger.WithPreffix("api.v1").Info("registered")
}
