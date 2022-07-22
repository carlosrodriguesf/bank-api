//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package transfer

import (
	"context"
	"github.com/carlosrodriguesf/bank-api/pkg/apputil/transaction"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
)

type (
	Options struct {
		Logger logger.Logger
		DB     db.Connection
	}
	Repository interface {
		Create(ctx context.Context, movement model.Transfer) (*model.GeneratedData, error)
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

func (r *repositoryImpl) WithTransaction(conn transaction.Transaction) Repository {
	return &repositoryImpl{
		logger: r.logger,
		db:     conn,
	}
}
