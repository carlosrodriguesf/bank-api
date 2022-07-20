//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package account

import (
	"context"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
)

type (
	Options struct {
		Logger logger.Logger
		DB     db.ExtendedDB
	}

	Repository interface {
		Create(ctx context.Context, account model.Account) (*model.GeneratedData, error)
		HasDocument(ctx context.Context, document string) (bool, error)
	}

	repositoryImpl struct {
		logger logger.Logger
		db     db.ExtendedDB
	}
)

func NewRepository(opts Options) Repository {
	return &repositoryImpl{
		logger: opts.Logger.WithLocation().WithPreffix("repository.account"),
		db:     opts.DB,
	}
}

func (r repositoryImpl) Create(ctx context.Context, account model.Account) (*model.GeneratedData, error) {
	generatedData := new(model.GeneratedData)
	query := `
		INSERT INTO accounts(name, document, secret, secret_salt) 
		VALUES (:name, :document, :secret, :secret_salt)
		RETURNING id, created_at`
	err := r.db.NamedGetContext(ctx, query, generatedData, account)
	if err != nil {
		r.logger.Error(err)
		return nil, err
	}
	return generatedData, nil
}

func (r repositoryImpl) HasDocument(ctx context.Context, document string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT TRUE FROM accounts WHERE document = $1)"
	err := r.db.GetContext(ctx, &exists, query, document)
	if err != nil {
		r.logger.Error(err)
	}
	return exists, err
}
