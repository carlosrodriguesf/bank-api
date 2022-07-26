//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package account

import (
	"context"
	"database/sql"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/transaction"
)

type (
	Options struct {
		Logger logger.Logger
		DB     db.Connection
	}

	Repository interface {
		Create(ctx context.Context, account model.Account) (*model.GeneratedData, error)
		HasDocument(ctx context.Context, document string) (bool, error)
		List(ctx context.Context) ([]model.Account, error)
		GetByIDOrDocument(ctx context.Context, v string) (*model.Account, error)
		UpdateBalance(ctx context.Context, accountID string, balance int64) error
		WithTransaction(conn transaction.Transaction) Repository
	}

	repositoryImpl struct {
		logger logger.Logger
		db     db.Connection
	}
)

func NewRepository(opts Options) Repository {
	return &repositoryImpl{
		logger: opts.Logger.WithLocation().WithPreffix("repository.account"),
		db:     opts.DB,
	}
}

func (r *repositoryImpl) Create(ctx context.Context, account model.Account) (*model.GeneratedData, error) {
	generatedData := new(model.GeneratedData)
	query := `
		INSERT INTO accounts(name, document, balance, secret, secret_salt) 
		VALUES (:name, :document, :balance, :secret, :secret_salt)
		RETURNING id, created_at`
	err := r.db.NamedGetContext(ctx, query, generatedData, account)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return generatedData, nil
}

func (r *repositoryImpl) HasDocument(ctx context.Context, document string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT TRUE FROM accounts WHERE document = $1)"
	err := r.db.GetContext(ctx, &exists, query, document)
	if err != nil {
		r.logger.Error(err)
	}
	return exists, err
}

func (r *repositoryImpl) List(ctx context.Context) ([]model.Account, error) {
	query := `SELECT id, name, document, balance, created_at FROM accounts`
	accounts := make([]model.Account, 0)
	err := r.db.SelectContext(ctx, &accounts, query)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return accounts, nil
}

func (r *repositoryImpl) GetByIDOrDocument(ctx context.Context, v string) (*model.Account, error) {
	query := "SELECT id, name, document, balance, secret, secret_salt, created_at FROM accounts WHERE id = $1 OR document = $1"
	acc := new(model.Account)
	err := r.db.GetContext(ctx, acc, query, v)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error(err)
		return nil, err
	}
	return acc, nil
}

func (r *repositoryImpl) UpdateBalance(ctx context.Context, accountID string, balance int64) error {
	query := "UPDATE accounts SET balance = $1 WHERE id = $2"
	_, err := r.db.ExecContext(ctx, query, balance, accountID)
	if err != nil {
		r.logger.Error(err)
	}
	return err
}

func (r *repositoryImpl) WithTransaction(conn transaction.Transaction) Repository {
	return &repositoryImpl{
		logger: r.logger,
		db:     conn,
	}
}
