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
		INSERT INTO accounts(name, document, balance, secret, secret_salt) 
		VALUES (?, ?, ?, ?, ?)
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
				Balance:    100,
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
					WithArgs("John Doe", "123.123.123-12", 100, "secret", "secret salt").
					WillReturnRows(rows)
			},
		},
		"should return error": {
			InputData: model.Account{
				Name:       "John Doe",
				Document:   "123.123.123-12",
				Balance:    100,
				Secret:     "secret",
				SecretSalt: "secret salt",
			},
			ExpectedData:  nil,
			ExpectedError: errors.New("fail"),
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs("John Doe", "123.123.123-12", 100, "secret", "secret salt").
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

func TestListAccounts(t *testing.T) {
	var (
		query           = regexp.QuoteMeta(`SELECT id, name, document, balance, created_at FROM accounts`)
		accountsExample = []model.Account{
			{
				ID:       "account_id_1",
				Name:     "Account Test 1",
				Document: "12312312312",
				Balance:  453,
			},
			{
				ID:       "account_id",
				Name:     "Account Test",
				Document: "12312312312",
				Balance:  819,
			},
		}
	)

	cases := map[string]struct {
		ExpectedData  []model.Account
		ExpectedError error
		PrepareMockDB func(mock sqlmock.Sqlmock)
	}{
		"should return success": {
			ExpectedData:  accountsExample,
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "document", "balance", "created_at"})
				for _, accountExample := range accountsExample {
					rows.AddRow(
						accountExample.ID,
						accountExample.Name,
						accountExample.Document,
						accountExample.Balance,
						accountExample.CreatedAt,
					)
				}
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
		},
		"should return success without accounts": {
			ExpectedData:  make([]model.Account, 0),
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "document", "balance", "created_at"})
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
		},
		"should return error": {
			ExpectedData:  nil,
			ExpectedError: errors.New("fail"),
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(query).WillReturnError(errors.New("fail"))
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

			data, err := repo.List(context.Background())

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}

func TestGetByIDOrDocument(t *testing.T) {
	var (
		query          = regexp.QuoteMeta(`SELECT id, name, document, balance, created_at FROM accounts WHERE id = $1 OR document = $1`)
		accountExample = model.Account{
			ID:         "account_id",
			Name:       "Account Test",
			Document:   "12312312312",
			Balance:    100,
			Secret:     "secret",
			SecretSalt: "secret_salt",
		}
	)

	cases := map[string]struct {
		InputData     string
		ExpectedData  *model.Account
		ExpectedError error
		PrepareMockDB func(mock sqlmock.Sqlmock)
	}{
		"should return success": {
			InputData:     "id_or_document",
			ExpectedData:  &accountExample,
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "name", "document", "balance", "secret", "secret_salt", "created_at"}).
					AddRow(
						accountExample.ID,
						accountExample.Name,
						accountExample.Document,
						accountExample.Balance,
						accountExample.Secret,
						accountExample.SecretSalt,
						accountExample.CreatedAt,
					)
				mock.
					ExpectQuery(query).
					WithArgs("id_or_document").
					WillReturnRows(rows)
			},
		},
		"should return success: account not found": {
			InputData:     "id_or_document",
			ExpectedData:  nil,
			ExpectedError: nil,
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "name", "document", "balance", "secret", "secret_salt", "created_at"})
				mock.
					ExpectQuery(query).
					WithArgs("id_or_document").
					WillReturnRows(rows)
			},
		},
		"should return error": {
			InputData:     "id_or_document",
			ExpectedData:  nil,
			ExpectedError: errors.New("fail"),
			PrepareMockDB: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectQuery(query).
					WithArgs("id_or_document").
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

			data, err := repo.GetByIDOrDocument(context.Background(), cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}
