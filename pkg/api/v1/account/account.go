package account

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/account"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	logger     logger.Logger
	accountApp account.App
}

func Register(g *echo.Group, opts apimodel.Options) {
	log := opts.Logger.WithPreffix("api.v1.account")
	h := handler{
		logger:     log.WithLocation(),
		accountApp: opts.App.Account(),
	}

	g.POST("/accounts", h.postAccount)
	g.GET("/accounts", h.getAccounts)
	g.GET("/accounts/:id/balance", h.getAccountBalance)

	log.Info("registered")
}

func (h *handler) postAccount(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)

	body := new(postAccountBody)
	if err := c.Bind(body); err != nil {
		log.Error(err)
		return apierror.ErrInvalidPayload
	}

	data, err := h.accountApp.Create(ctx, model.Account{
		Name:     body.Name,
		Document: body.Document,
		Secret:   body.Secret,
		Balance:  body.Balance,
	})
	if err != nil {
		if err := apierror.Get(err, errorMap); err != nil {
			return err
		}
		log.Error(err)
		return apierror.ErrInternal
	}
	return c.JSON(http.StatusOK, apimodel.Response{
		Data: data,
	})
}

func (h *handler) getAccounts(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)

	data, err := h.accountApp.List(ctx)
	if err != nil {
		if err := apierror.Get(err, errorMap); err != nil {
			return err
		}
		log.Error(err)
		return apierror.ErrInternal
	}
	return c.JSON(http.StatusOK, apimodel.Response{
		Data: data,
	})
}

func (h *handler) getAccountBalance(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)
	accountID := c.Param("id")
	data, err := h.accountApp.GetBalance(ctx, accountID)
	if err != nil {
		if err := apierror.Get(err, errorMap); err != nil {
			return err
		}
		log.Error(err)
		return apierror.ErrInternal
	}
	return c.JSON(http.StatusOK, apimodel.Response{
		Data: data,
	})
}
