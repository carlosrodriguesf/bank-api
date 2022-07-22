package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	apierror "github.com/carlosrodriguesf/bank-api/pkg/api/error"
	apimodel "github.com/carlosrodriguesf/bank-api/pkg/api/model"
	"github.com/carlosrodriguesf/bank-api/pkg/app/auth"
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

func TestHandler_Login(t *testing.T) {
	var (
		endpoint           = "/api/v1/login"
		credentialsExample = model.Credentials{
			Document: "123.123.123-18",
			Secret:   "1234",
		}
		accountExample = model.Account{
			ID:       "account_id",
			Name:     "Test Account",
			Document: "12312312312",
		}
		sessionExample = model.Session{
			Token:     "generated_token",
			Account:   accountExample,
			CreatedAt: time.Now(),
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
		ExpectedData   *model.Session
		ExpectedErr    error
		PrepareMockApp func(mock *auth.MockApp)
	}{
		"should return success": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(credentialsExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: &sessionExample,
			ExpectedErr:  nil,
			PrepareMockApp: func(mock *auth.MockApp) {
				mock.EXPECT().
					Auth(gomock.Any(), credentialsExample).
					Return(&sessionExample, nil)
			},
		},
		"should return error on bind": {
			InputData: func(t *testing.T) io.Reader {
				return strings.NewReader("invalid body")
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInvalidPayload,
			PrepareMockApp: func(mock *auth.MockApp) {
			},
		},
		"should return validation error": {
			InputData: func(t *testing.T) io.Reader {
				credentialsExample := credentialsExample
				credentialsExample.Document = ""
				body, err := json.Marshal(credentialsExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.NewApiError(http.StatusBadRequest, "invalid_payload", validationErrorExample.Violations),
			PrepareMockApp: func(mock *auth.MockApp) {
				credentialsExample := credentialsExample
				credentialsExample.Document = ""
				mock.EXPECT().
					Auth(gomock.Any(), credentialsExample).
					Return(nil, validationErrorExample)
			},
		},
		"should return error: invalid credentials": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(credentialsExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  errorMap[pkgerror.ErrInvalidCredentials],
			PrepareMockApp: func(mock *auth.MockApp) {
				mock.EXPECT().
					Auth(gomock.Any(), credentialsExample).
					Return(nil, pkgerror.ErrInvalidCredentials)
			},
		},
		"should return error": {
			InputData: func(t *testing.T) io.Reader {
				body, err := json.Marshal(credentialsExample)
				assert.NoError(t, err)
				return bytes.NewReader(body)
			},
			ExpectedData: nil,
			ExpectedErr:  apierror.ErrInternal,
			PrepareMockApp: func(mock *auth.MockApp) {
				mock.EXPECT().
					Auth(gomock.Any(), credentialsExample).
					Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl, ctx := gomock.WithContext(context.Background(), t)

			mockApp := auth.NewMockApp(ctrl)

			cs.PrepareMockApp(mockApp)

			h := handler{
				logger:  logger.New(""),
				authApp: mockApp,
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, endpoint, cs.InputData(t)).WithContext(ctx)
			rec := httptest.NewRecorder()
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			c := e.NewContext(req, rec)
			c.SetPath(endpoint)

			err := h.login(c)

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
