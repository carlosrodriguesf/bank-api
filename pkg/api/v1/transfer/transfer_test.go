package transfer

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/transfer"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandler_CreateTransfer(t *testing.T) {
	var (
		endpoint            = "/api/v1/transfers"
		postTransferExample = postTransferBody{
			TargetAccountID: "target_account_id",
			Amount:          500,
		}
		createTransferExample = model.Transfer{
			OriginAccountID: "origin_account_id",
			TargetAccountID: postTransferExample.TargetAccountID,
			Amount:          postTransferExample.Amount,
		}
		createdTransferExample = model.Transfer{
			ID:              "transfer_id",
			OriginAccountID: createTransferExample.OriginAccountID,
			TargetAccountID: createTransferExample.TargetAccountID,
		}
		validationErrorExample = &validator.ValidationError{
			OriginalMessage: "invalid data",
			Message:         "invalid data",
			Violations: []validator.Violation{{
				Field: "name",
				Tag:   "required",
			}},
		}
	)

	cases := map[string]struct {
		InputData      func(t *testing.T) io.Reader
		ExpectedData   *model.Transfer
		ExpectedErr    error
		PrepareMockApp func(mock *transfer.MockApp)
	}{
		"should return success": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: &createdTransferExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(&createdTransferExample, nil)
			},
		},
		"should return error on bind": {
			InputData: func(t *testing.T) io.Reader {
				return strings.NewReader("invalid body")
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInvalidPayload,
			PrepareMockApp: func(mock *transfer.MockApp) {
			},
		},
		"should return validation error": {
			InputData: func(t *testing.T) io.Reader {
				postTransferExample := postTransferExample
				postTransferExample.TargetAccountID = ""
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.NewApiError(http.StatusBadRequest, "invalid_payload", validationErrorExample.Violations),
			PrepareMockApp: func(mock *transfer.MockApp) {
				createTransferExample := createTransferExample
				createTransferExample.TargetAccountID = ""
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(nil, validationErrorExample)
			},
		},
		"should return error: origin account not found": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrOriginAccountTransferNotFound],
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(nil, pkgerror.ErrOriginAccountTransferNotFound)
			},
		},
		"should return error: target account not found": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrOriginAccountTransferNotFound],
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(nil, pkgerror.ErrOriginAccountTransferNotFound)
			},
		},
		"should return error: unsifficient funds": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrInsufficientFunds],
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(nil, pkgerror.ErrInsufficientFunds)
			},
		},
		"should return error": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postTransferExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createTransferExample).
					Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			ctx = model.SetSessionOnContext(ctx, &model.Session{
				Token: "session_token",
				Account: model.Account{
					ID: createTransferExample.OriginAccountID,
				},
			})

			mockApp := transfer.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:      logger.New(""),
				transferApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, cs.InputData(t)).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := h.postTransfer(c)

			assert.Equal(t, cs.ExpectedErr, err)

			expectedResponseJSON, err := json.Marshal(apimodel.Response{Data: cs.ExpectedData})
			assert.NoError(t, err)

			var expectedResponse apimodel.Response
			err = json.Unmarshal(expectedResponseJSON, &expectedResponse)
			assert.NoError(t, err)

			var currentResponse apimodel.Response
			json.NewDecoder(rec.Body).Decode(&currentResponse)

			assert.Equal(t, expectedResponse, currentResponse)
		})
	}
}

func TestHandler_getTransfers(t *testing.T) {
	var (
		endpoint         = "/api/v1/transfers"
		transfersExample = []model.TransferDetailed{{
			Transfer: model.Transfer{
				ID:              "transfer_id",
				OriginAccountID: "origin_account_id",
				TargetAccountID: "target_account_id",
				Amount:          500,
				CreatedAt:       time.Now(),
			},
			OriginAccountName: "Origin Account",
			TargetAccountName: "Target account",
		}}
	)

	cases := map[string]struct {
		ExpectedData   []model.TransferDetailed
		ExpectedErr    error
		PrepareMockApp func(mock *transfer.MockApp)
	}{
		"should return success": {
			ExpectedData: transfersExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().List(gomock.Any(), "origin_account_id").Return(transfersExample, nil)
			},
		},
		"should return error": {
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrCantListTransfers],
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().List(gomock.Any(), "origin_account_id").Return(nil, pkgerror.ErrCantListTransfers)
			},
		},
		"should return internal error": {
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *transfer.MockApp) {
				mock.EXPECT().List(gomock.Any(), "origin_account_id").Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			ctx = model.SetSessionOnContext(ctx, &model.Session{
				Token: "session_token",
				Account: model.Account{
					ID: "origin_account_id",
				},
			})

			mockApp := transfer.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:      logger.New(""),
				transferApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := h.getTransfers(c)

			assert.Equal(t, cs.ExpectedErr, err)

			expectedResponseJSON, err := json.Marshal(apimodel.Response{Data: cs.ExpectedData})
			assert.NoError(t, err)

			var expectedResponse apimodel.Response
			err = json.Unmarshal(expectedResponseJSON, &expectedResponse)
			assert.NoError(t, err)

			var currentResponse apimodel.Response
			json.NewDecoder(rec.Body).Decode(&currentResponse)

			assert.Equal(t, expectedResponse, currentResponse)
		})
	}
}
