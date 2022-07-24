//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package transfer

import (
	"context"
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
		Create(ctx context.Context, movement model.Transfer) (*model.GeneratedData, error)
		List(ctx context.Context, accountID string) ([]model.TransferDetailed, error)
		WithTransaction(conn transaction.Transaction) Repository
	}
	repositoryImpl struct {
		logger logger.Logger
		db     db.Connection
	}
)

func NewRepository(opts Options) Repository {
	return &repositoryImpl{
		logger: opts.Logger,
		db:     opts.DB,
	}
}

func (r *repositoryImpl) Create(ctx context.Context, movement model.Transfer) (*model.GeneratedData, error) {
	query := `
		INSERT INTO transfers(origin_account_id, target_account_id, amount) 
		VALUES (:origin_account_id, :target_account_id, :amount)
		RETURNING id, created_at`
	generatedData := new(model.GeneratedData)
	err := r.db.NamedGetContext(ctx, query, generatedData, movement)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return generatedData, nil
}

func (r *repositoryImpl) List(ctx context.Context, accountID string) ([]model.TransferDetailed, error) {
	query := `
		SELECT 
			t.id, 
			t.origin_account_id, 
			t.target_account_id, 
			t.amount, 
			t.created_at, 
			t.origin_account_id = $1 AS sent,
			oa.name AS origin_account_name,
			ta.name AS target_account_name
		FROM transfers t
			INNER JOIN accounts oa ON oa.id = t.origin_account_id
			INNER JOIN accounts ta ON ta.id = t.target_account_id
		WHERE origin_account_id = $1 OR target_account_id = $1`
	transfers := make([]model.TransferDetailed, 0)
	err := r.db.SelectContext(ctx, &transfers, query, accountID)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return transfers, nil
}

func (r *repositoryImpl) WithTransaction(conn transaction.Transaction) Repository {
	return &repositoryImpl{
		logger: r.logger,
		db:     conn,
	}
}
