package account

import (
	"context"
	"errors"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	var (
		currentTime         = time.Now()
		creationDataExample = model.Account{
			Name:     "John Doe",
			Document: "123.123.123-12",
			Secret:   "secret",
		}
		accountExample = model.Account{
			ID:         "account_id",
			Name:       "John Doe",
			Document:   "12312312312",
			Secret:     "generated_secret",
			SecretSalt: "generated_salt",
			CreatedAt:  currentTime,
		}
		validationErrorExample = &validator.ValidationError{
			OriginalMessage: "fail",
			Message:         "fail",
			Violations: []validator.Violation{{
				Field: "Secret",
				Tag:   "required",
			}},
		}
	)

	cases := map[string]struct {
		InputData             model.Account
		ExpectedData          *model.Account
		ExpectedError         error
		PrepareMockValidator  func(mock *validator.MockValidator)
		PrepareMockSecret     func(mock *secret.MockSecret)
		PrepareMockRepository func(mock *account.MockRepository)
	}{
		"should return success": {
			InputData:     creationDataExample,
			ExpectedData:  &accountExample,
			ExpectedError: nil,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(creationDataExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().GenSalt().Return(accountExample.SecretSalt)
				mock.EXPECT().
					Encode(creationDataExample.Secret, accountExample.SecretSalt).
					Return(accountExample.Secret)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				mock.EXPECT().
					HasDocument(gomock.Any(), model.DocumentRegex.ReplaceAllString(creationDataExample.Document, "")).
					Return(false, nil)
				mock.EXPECT().
					Create(gomock.Any(), model.Account{
						Name:       accountExample.Name,
						Document:   accountExample.Document,
						Secret:     accountExample.Secret,
						SecretSalt: accountExample.SecretSalt,
					}).
					Return(&model.GeneratedData{
						ID:        accountExample.ID,
						CreatedAt: accountExample.CreatedAt,
					}, nil)
			},
		},
		"should return error on validate": {
			InputData:     creationDataExample,
			ExpectedData:  nil,
			ExpectedError: validationErrorExample,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(creationDataExample).Return(validationErrorExample)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {

			},
			PrepareMockRepository: func(mock *account.MockRepository) {

			},
		},
		"should return error on save account": {
			InputData:     creationDataExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateAccount,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(creationDataExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().GenSalt().Return(accountExample.SecretSalt)
				mock.EXPECT().
					Encode(creationDataExample.Secret, accountExample.SecretSalt).
					Return(accountExample.Secret)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				mock.EXPECT().
					HasDocument(gomock.Any(), model.DocumentRegex.ReplaceAllString(creationDataExample.Document, "")).
					Return(false, nil)
				mock.EXPECT().
					Create(gomock.Any(), model.Account{
						Name:       accountExample.Name,
						Document:   accountExample.Document,
						Secret:     accountExample.Secret,
						SecretSalt: accountExample.SecretSalt,
					}).
					Return(nil, errors.New("fail"))
			},
		},
		"should return error on check document": {
			InputData:     creationDataExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateAccount,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(creationDataExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().GenSalt().Return(accountExample.SecretSalt)
				mock.EXPECT().
					Encode(creationDataExample.Secret, accountExample.SecretSalt).
					Return(accountExample.Secret)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				mock.EXPECT().
					HasDocument(gomock.Any(), model.DocumentRegex.ReplaceAllString(creationDataExample.Document, "")).
					Return(false, errors.New("fail"))
			},
		},
		"should return error: document exists": {
			InputData:     creationDataExample,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrDocumentAlreadyExists,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(creationDataExample).Return(nil)
			},
			PrepareMockSecret: func(mock *secret.MockSecret) {
				mock.EXPECT().GenSalt().Return(accountExample.SecretSalt)
				mock.EXPECT().
					Encode(creationDataExample.Secret, accountExample.SecretSalt).
					Return(accountExample.Secret)
			},
			PrepareMockRepository: func(mock *account.MockRepository) {
				mock.EXPECT().
					HasDocument(gomock.Any(), model.DocumentRegex.ReplaceAllString(creationDataExample.Document, "")).
					Return(true, nil)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl, ctx      = gomock.WithContext(context.Background(), t)
				mockSecret     = secret.NewMockSecret(ctrl)
				mockValidator  = validator.NewMockValidator(ctrl)
				mockRepository = account.NewMockRepository(ctrl)
			)

			cs.PrepareMockSecret(mockSecret)
			cs.PrepareMockValidator(mockValidator)
			cs.PrepareMockRepository(mockRepository)

			service := NewApp(Options{
				Logger:      logger.New(""),
				Secret:      mockSecret,
				Validator:   mockValidator,
				RepoAccount: mockRepository,
			})

			data, err := service.Create(ctx, cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}
