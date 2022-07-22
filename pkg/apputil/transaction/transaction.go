//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package transaction

import (
	"context"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
)

type (
	Transaction = db.ExtendedTx
	Manager     interface {
		Create(ctx context.Context) (Transaction, error)
		EndTransaction(tx Transaction, commit bool) error
	}
	manager struct {
		db db.ExtendedDB
	}
)

func NewManager(db db.ExtendedDB) Manager {
	return &manager{
		db: db,
	}
}

func (r *manager) Create(ctx context.Context) (Transaction, error) {
	return db.BeginTransaction(ctx, r.db)
}

func (r *manager) EndTransaction(tx Transaction, commit bool) error {
	if commit {
		return tx.Commit()
	}
	return tx.Rollback()
}
