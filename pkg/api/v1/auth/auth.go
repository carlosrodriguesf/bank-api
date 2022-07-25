package auth

import (
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/auth"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

type handler struct {
	logger  logger.Logger
	authApp auth.App
}

func Register(g *echo.Group, opts apimodel.Options) {
	log := opts.Logger.WithPreffix("api.v1.account")

	h := handler{
		logger:  log.WithLocation(),
		authApp: opts.App.Auth(),
	}

	g.POST("/login", h.login)

	log.Info("registered")
}

// login swagger document
// @Description Login
// @Tags auth
// @Produce json
// @Param credentials body model.Credentials true "expected structure"
// @Success 200 {object} model.Response{data=model.Account}
// @Success 400 {object} model.Response{error=error.ApiError}
// @Failure 500 {object} model.Response{error=error.ApiError}
// @Router /api/v1/login [post]
func (h *handler) login(c echo.Context) error {
	ctx := c.Request().Context()
	log := h.logger.WithContext(ctx)

	var body model.Credentials
	if err := c.Bind(&body); err != nil {
		log.Error(err)
		return apierror.ErrInvalidPayload
	}

	data, err := h.authApp.Auth(ctx, body)
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
