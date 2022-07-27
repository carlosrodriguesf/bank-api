package auth

import (
	"context"
	"errors"
	"fmt"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/cache"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/generate"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAuth(t *testing.T) {
	var (
		currentTime        = time.Now()
		uuidExample        = uuid.NewString()
		credentialsExample = model.Credentials{
			Document: "123.123.123-12",
			Secret:   "secret",
		}
		accountExample = model.Account{
			ID:        "account_id",
			Name:      "John Doe",
			Document:  "123.123.123-12",
			CreatedAt: currentTime,
		}
		sessionExample = model.Session{
			Token: uuidExample,
			Account: model.Account{
				ID:        accountExample.ID,
				Name:      accountExample.Name,
				Document:  accountExample.Document,
				CreatedAt: accountExample.CreatedAt,
			},
			CreatedAt: currentTime,
		}
	)

	cases := map[string]struct {
		InputData             model.Credentials
		ExpectedData          *model.Session
		ExpectedError         error
		PrepareMockValidator  func(mock *validator.MockValidator)
		PrepareMockSecret     func(mock *secret.MockSecret)
		PrepareMockRepository func(mock *account.MockRepository)
		PrepareMockCache      func(mock *cache.MockCache)
	}{
		"should return success": {
			InputData:     credentialsExample,
			ExpectedData:  &sessionExample,
			ExpectedError: nil,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().
					Verify(credentialsExample.Secret, accountExample.Secret, accountExample.SecretSalt).
					Return(true)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				document := model.DocumentRegex.ReplaceAllString(accountExample.Document, "")
				mock.EXPECT().
					GetByIDOrDocument(gomock.Any(), document).
					Return(&accountExample, nil)
			},
			PrepareMockCache: func(mock *cache.MockCache) {
				mock.EXPECT().
					Set(gomock.Any(), getSessionCacheKey(sessionExample.Token), &sessionExample, cacheExpiration).
					Return(nil)
			},
		},
		"should return validation error": {
			InputData:     credentialsExample,
			ExpectedData:  nil,
			ExpectedError: &validator.ValidationError{OriginalMessage: "fail"},
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(&validator.ValidationError{OriginalMessage: "fail"})
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
			},
			PrepareMockCache: func(mock *cache.MockCache) {
			},
		},
		"should return error on get account": {
			InputData:     credentialsExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantAuth,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				document := model.DocumentRegex.ReplaceAllString(accountExample.Document, "")
				mock.EXPECT().
					GetByIDOrDocument(gomock.Any(), document).
					Return(nil, errors.New("fail"))
			},
			PrepareMockCache: func(mock *cache.MockCache) {
			},
		},
		"should return error: document invalid": {
			InputData:     credentialsExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrInvalidCredentials,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				document := model.DocumentRegex.ReplaceAllString(accountExample.Document, "")
				mock.EXPECT().
					GetByIDOrDocument(gomock.Any(), document).
					Return(nil, nil)
			},
			PrepareMockCache: func(mock *cache.MockCache) {
			},
		},
		"should return error: password invalid": {
			InputData:     credentialsExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrInvalidCredentials,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().
					Verify(credentialsExample.Secret, accountExample.Secret, accountExample.SecretSalt).
					Return(false)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				document := model.DocumentRegex.ReplaceAllString(accountExample.Document, "")
				mock.EXPECT().
					GetByIDOrDocument(gomock.Any(), document).
					Return(&accountExample, nil)
			},
			PrepareMockCache: func(mock *cache.MockCache) {
			},
		},
		"should return error on save session": {
			InputData:     credentialsExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantAuth,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(credentialsExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().
					Verify(credentialsExample.Secret, accountExample.Secret, accountExample.SecretSalt).
					Return(true)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				document := model.DocumentRegex.ReplaceAllString(accountExample.Document, "")
				mock.EXPECT().
					GetByIDOrDocument(gomock.Any(), document).
					Return(&accountExample, nil)
			},
			PrepareMockCache: func(mock *cache.MockCache) {
				mock.EXPECT().
					Set(gomock.Any(), fmt.Sprintf(cacheKeySession, sessionExample.Token), &sessionExample, cacheExpiration).
					Return(errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl, ctx      = gomock.WithContext(context.Background(), t)
				mockCache      = cache.NewMockCache(ctrl)
				mockSecret     = secret.NewMockSecret(ctrl)
				mockValidator  = validator.NewMockValidator(ctrl)
				mockRepository = account.NewMockRepository(ctrl)
				mockGenerate   = generate.NewMockGenerate(ctrl)
			)

			cs.PrepareMockSecret(mockSecret)
			cs.PrepareMockValidator(mockValidator)
			cs.PrepareMockRepository(mockRepository)
			cs.PrepareMockCache(mockCache)

			mockGenerate.EXPECT().UUID().AnyTimes().Return(uuidExample)
			mockGenerate.EXPECT().CurrentTime().AnyTimes().Return(currentTime)

			app := NewApp(Options{
				Logger:      logger.New(""),
				Secret:      mockSecret,
				Cache:       mockCache,
				Validator:   mockValidator,
				RepoAccount: mockRepository,
				Generate:    mockGenerate,
			})

			data, err := app.Auth(ctx, cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}

func TestGetSessionByToken(t *testing.T) {
	var (
		currentTime    = time.Now()
		uuidExample    = uuid.NewString()
		sessionExample = model.Session{
			Token: uuidExample,
			Account: model.Account{
				ID:        "account_id",
				Name:      "John Doe",
				Document:  "123.123.123-12",
				CreatedAt: currentTime,
			},
			CreatedAt: currentTime,
		}
	)

	cases := map[string]struct {
		InputData        string
		ExpectedData     *model.Session
		ExpectedError    error
		PrepareMockCache func(mock *cache.MockCache)
	}{
		"should return success": {
			InputData:     uuidExample,
			ExpectedData:  &sessionExample,
			ExpectedError: nil,
			PrepareMockCache: func(mock *cache.MockCache) {
				mock.EXPECT().
					GetUpdating(gomock.Any(), getSessionCacheKey(uuidExample), new(model.Session), cacheExpiration).
					Do(func(_, _ interface{}, session *model.Session, _ time.Duration) {
						*session = sessionExample
					})
			},
		},
		"should return cache missing error": {
			InputData:     uuidExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrSessionNotFound,
			PrepareMockCache: func(mock *cache.MockCache) {
				mock.EXPECT().
					GetUpdating(gomock.Any(), getSessionCacheKey(uuidExample), new(model.Session), cacheExpiration).
					Return(errors.New("cache missing"))
				mock.EXPECT().
					IsErrCacheMissing(errors.New("cache missing")).
					Return(true)
			},
		},
		"should return cant get session error": {
			InputData:     uuidExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantGetSession,
			PrepareMockCache: func(mock *cache.MockCache) {
				mock.EXPECT().
					GetUpdating(gomock.Any(), getSessionCacheKey(uuidExample), new(model.Session), cacheExpiration).
					Return(errors.New("cache missing"))
				mock.EXPECT().
					IsErrCacheMissing(errors.New("cache missing")).
					Return(false)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl, ctx      = gomock.WithContext(context.Background(), t)
				mockCache      = cache.NewMockCache(ctrl)
				mockSecret     = secret.NewMockSecret(ctrl)
				mockValidator  = validator.NewMockValidator(ctrl)
				mockRepository = account.NewMockRepository(ctrl)
				mockGenerate   = generate.NewMockGenerate(ctrl)
			)

			cs.PrepareMockCache(mockCache)

			mockGenerate.EXPECT().UUID().AnyTimes().Return(uuidExample)
			mockGenerate.EXPECT().CurrentTime().AnyTimes().Return(currentTime)

			app := NewApp(Options{
				Logger:      logger.New(""),
				Secret:      mockSecret,
				Cache:       mockCache,
				Validator:   mockValidator,
				RepoAccount: mockRepository,
				Generate:    mockGenerate,
			})

			data, err := app.GetSessionByToken(ctx, cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}
