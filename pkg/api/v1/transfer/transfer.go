package transfer

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/transfer"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	logger      logger.Logger
	transferApp transfer.App
}

func Register(g *echo.Group, opts apimodel.Options) {
	log := opts.Logger.WithPreffix("api.v1.transfer")
	h := handler{
		logger:      opts.Logger.WithLocation(),
		transferApp: opts.App.Transfer(),
	}

	g.POST("/transfers", h.postTransfer, opts.Middleware.Auth().Private)
	g.GET("/transfers", h.getTransfers, opts.Middleware.Auth().Private)

	log.Info("registered")
}

// postTransfer swagger document
// @Description Make a transfer between accounts
// @Tags transfer
// @Produce json
// @Security UserToken
// @Param transfer body postTransferBody true "expected structure"
// @Success 200 {object} model.Response{data=model.Transfer}
// @Success 400 {object} model.Response{error=error.ApiError}
// @Failure 500 {object} model.Response{error=error.ApiError}
// @Router /api/v1/transfers [post]
func (h *handler) postTransfer(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)

	body := new(postTransferBody)
	if err := c.Bind(body); err != nil {
		log.Error(err)
		return apierror.ErrInvalidPayload
	}

	sess := model.GetSessionFromContext(ctx)
	data, err := h.transferApp.Create(ctx, model.Transfer{
		OriginAccountID: sess.Account.ID,
		TargetAccountID: body.TargetAccountID,
		Amount:          body.Amount,
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

// getTransfers swagger document
// @Description List of transfer received or sent by current auth user
// @Tags transfer
// @Produce json
// @Security UserToken
// @Success 200 {object} model.Response{data=[]model.Account}
// @Failure 500 {object} model.Response{error=error.ApiError}
// @Router /api/v1/transfers [get]
func (h *handler) getTransfers(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)

	sess := model.GetSessionFromContext(ctx)
	data, err := h.transferApp.List(ctx, sess.Account.ID)
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
