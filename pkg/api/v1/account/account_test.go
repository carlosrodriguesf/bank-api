package account

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/account"
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
)

func TestHandler_CreateAccount(t *testing.T) {
	var (
		endpoint           = "/api/v1/accounts"
		postAccountExample = postAccountBody{
			Name:     "John Doe",
			Document: "123.123.123-18",
			Secret:   "1234",
		}
		createAccountExample = model.Account{
			Name:     postAccountExample.Name,
			Document: postAccountExample.Document,
			Secret:   postAccountExample.Secret,
		}
		createdAccountExample = model.Account{
			ID:         "generated_id",
			Name:       createAccountExample.Name,
			Document:   createAccountExample.Document,
			Secret:     "generated_secret",
			SecretSalt: "generated_secret_salt",
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
		ExpectedData   *model.Account
		ExpectedErr    error
		PrepareMockApp func(mock *account.MockApp)
	}{
		"should return success": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postAccountExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: &createdAccountExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createAccountExample).
					Return(&createdAccountExample, nil)
			},
		},
		"should return error on bind": {
			InputData: func(t *testing.T) io.Reader {
				return strings.NewReader("invalid body")
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInvalidPayload,
			PrepareMockApp: func(mock *account.MockApp) {
			},
		},
		"should return validation error": {
			InputData: func(t *testing.T) io.Reader {
				postAccountExample := postAccountExample
				postAccountExample.Name = ""
				body, err := json.Marshal(postAccountExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.NewApiError(http.StatusBadRequest, "invalid_payload", validationErrorExample.Violations),
			PrepareMockApp: func(mock *account.MockApp) {
				createAccountExample := createAccountExample
				createAccountExample.Name = ""
				mock.EXPECT().
					Create(gomock.Any(), createAccountExample).
					Return(nil, validationErrorExample)
			},
		},
		"should return error: document already exists": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postAccountExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrDocumentAlreadyExists],
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createAccountExample).
					Return(nil, pkgerror.ErrDocumentAlreadyExists)
			},
		},
		"should return error": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(postAccountExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					Create(gomock.Any(), createAccountExample).
					Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockApp := account.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:     logger.New(""),
				accountApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, cs.InputData(t)).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := h.postAccount(c)

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

func TestHandler_getAccounts(t *testing.T) {
	var (
		endpoint        = "/api/v1/accounts"
		accountsExample = []model.Account{
			{
				ID:       "account_id_1",
				Name:     "Account Test 1",
				Document: "12312312312",
			},
			{
				ID:       "account_id",
				Name:     "Account Test",
				Document: "12312312312",
			},
		}
	)

	cases := map[string]struct {
		ExpectedData   []model.Account
		ExpectedErr    error
		PrepareMockApp func(mock *account.MockApp)
	}{
		"should return success": {
			ExpectedData: accountsExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().List(gomock.Any()).Return(accountsExample, nil)
			},
		},
		"should return error": {
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrCantListAccounts],
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().List(gomock.Any()).Return(nil, pkgerror.ErrCantListAccounts)
			},
		},
		"should return internal error": {
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().List(gomock.Any()).Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockApp := account.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:     logger.New(""),
				accountApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := h.getAccounts(c)

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

func TestHandler_getAccountBalance(t *testing.T) {
	var (
		endpoint       = "/api/v1/accounts/account_id/balance"
		accountID      = "account_id"
		balanceExample = model.AccountBalance{
			Balance: 1000,
		}
	)

	cases := map[string]struct {
		ExpectedData   *model.AccountBalance
		ExpectedErr    error
		PrepareMockApp func(mock *account.MockApp)
	}{
		"should return success": {
			ExpectedData: &balanceExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					GetBalance(gomock.Any(), accountID).
					Return(&balanceExample, nil)
			},
		},
		"should return error: account not found": {
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrAccountNotFound],
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					GetBalance(gomock.Any(), accountID).
					Return(nil, pkgerror.ErrAccountNotFound)
			},
		},
		"should return error: cant get account balance": {
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrCantGetAccountBalance],
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					GetBalance(gomock.Any(), accountID).
					Return(nil, pkgerror.ErrCantGetAccountBalance)
			},
		},
		"should return internal error": {
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *account.MockApp) {
				mock.EXPECT().
					GetBalance(gomock.Any(), accountID).
					Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockApp := account.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:     logger.New(""),
				accountApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, nil).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)
			c.SetParamNames("id")
			c.SetParamValues(accountID)

			err := h.getAccountBalance(c)

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
