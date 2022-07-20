package account

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/test"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	query := regexp.QuoteMeta(`
		INSERT INTO accounts(name, document, secret, secret_salt) 
		VALUES (?, ?, ?, ?)
		RETURNING id, created_at`)

	currentTime := time.Now()

	cases := map[string]struct {
		InputData     model.Account
		ExpectedData  *model.GeneratedData
		ExpectedError error
		PrepareMockDB func(mock sqlmock.Sqlmock)
	}{
		"should return success": {
			InputData: model.Account{
				Name:       "John Doe",
				Document:   "123.123.123-12",
				Secret:     "secret",
				SecretSalt: "secret salt",
			},
			ExpectedData: &model.GeneratedData{
				ID:        "generated_id",
				CreatedAt: currentTime,
			},
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "created_at"}).
					AddRow("generated_id", currentTime)
				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs("John Doe", "123.123.123-12", "secret", "secret salt").
					WillReturnRows(rows)
			},
		},
		"should return error": {
			InputData: model.Account{
				Name:       "John Doe",
				Document:   "123.123.123-12",
				Secret:     "secret",
				SecretSalt: "secret salt",
			},
			ExpectedData:  nil,
			ExpectedError: errors.New("fail"),
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs("John Doe", "123.123.123-12", "secret", "secret salt").
					WillReturnError(errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			conn, mock := test.GetSQLMock()

			cs.PrepareMockDB(mock)

			repo := NewRepository(Options{
				Logger: logger.New(""),
				DB:     db.NewExtendedDB(conn),
			})

			data, err := repo.Create(context.Background(), cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}

func TestHasDocument(t *testing.T) {
	query := regexp.QuoteMeta(`SELECT EXISTS(SELECT TRUE FROM accounts WHERE document = $1)`)

	cases := map[string]struct {
		InputData     string
		ExpectedData  bool
		ExpectedError error
		PrepareMockDB func(mock sqlmock.Sqlmock)
	}{
		"should return success": {
			InputData:     "12312312312",
			ExpectedData:  true,
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"exists"}).
					AddRow(true)
				mock.
					ExpectQuery(query).
					WithArgs("12312312312").
					WillReturnRows(rows)
			},
		},
		"should return error": {
			InputData:     "12312312312",
			ExpectedData:  false,
			ExpectedError: errors.New("fail"),
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(query).
					WithArgs("12312312312").
					WillReturnError(errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			conn, mock := test.GetSQLMock()

			cs.PrepareMockDB(mock)

			repo := NewRepository(Options{
				Logger: logger.New(""),
				DB:     db.NewExtendedDB(conn),
			})

			data, err := repo.HasDocument(context.Background(), cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}
