package v1

import (
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/api/v1/account"
	"github.com/carlosrodriguesf/bank-api/pkg/api/v1/auth"
	"github.com/carlosrodriguesf/bank-api/pkg/api/v1/transfer"
	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group, opts apimodel.Options) {
	log := opts.Logger.WithPreffix("api.v1")

	g = g.Group("/v1")

	account.Register(g, opts)
	auth.Register(g, opts)
	transfer.Register(g, opts)

	log.Info("registered")
}
