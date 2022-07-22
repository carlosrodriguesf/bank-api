package transfer

import (
	"context"
	"database/sql"
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
	var (
		currentTime     = time.Now()
		transferExample = model.Transfer{
			OriginAccountID: "origin_account_id",
			TargetAccountID: "target_account_id",
			Amount:          500,
		}
		genreatedDataExample = model.GeneratedData{
			ID:        "generated_id",
			CreatedAt: currentTime,
		}
		query = regexp.QuoteMeta(`
			INSERT INTO transfers(origin_account_id, target_account_id, amount) 
			VALUES (?, ?, ?)
			RETURNING id, created_at
		`)
	)
	cases := map[string]struct {
		InputData      model.Transfer
		ExpectedData   *model.GeneratedData
		ExpectedError  error
		PrepareMockSQL func(mock sqlmock.Sqlmock)
	}{
		"should return success": {
			InputData:     transferExample,
			ExpectedData:  &genreatedDataExample,
			ExpectedError: nil,
			PrepareMockSQL: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "created_at"}).
					AddRow(
						genreatedDataExample.ID,
						genreatedDataExample.CreatedAt,
					)
				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs(transferExample.OriginAccountID, transferExample.TargetAccountID, transferExample.Amount).
					WillReturnRows(rows)
			},
		},
		"should return no data error": {
			InputData:     transferExample,
			ExpectedData:  nil,
			ExpectedError: sql.ErrNoRows,
			PrepareMockSQL: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.
					NewRows([]string{"id", "created_at"})

				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs(transferExample.OriginAccountID, transferExample.TargetAccountID, transferExample.Amount).
					WillReturnRows(rows)
			},
		},
		"should return error": {
			InputData:     transferExample,
			ExpectedData:  nil,
			ExpectedError: errors.New("fail"),
			PrepareMockSQL: func(mock sqlmock.Sqlmock) {
				mock.ExpectPrepare(query).
					ExpectQuery().
					WithArgs(transferExample.OriginAccountID, transferExample.TargetAccountID, transferExample.Amount).
					WillReturnError(errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				dbConn, sqlMock = test.GetSQLMock()
				repository      = NewRepository(Options{
					Logger: logger.New(""),
					DB:     db.NewExtendedDB(dbConn),
				})
			)

			cs.PrepareMockSQL(sqlMock)

			data, err := repository.Create(context.Background(), transferExample)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}

func TestWithTransaction(t *testing.T) {
	repoWithDB := &repositoryImpl{
		db: db.ExtendedDB(nil),
	}
	repoWithTx := &repositoryImpl{
		db: db.ExtendedTx(nil),
	}
	assert.Equal(t, repoWithTx, repoWithDB.WithTransaction(db.ExtendedTx(nil)))
}
