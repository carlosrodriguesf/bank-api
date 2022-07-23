package repository

import (
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/transfer"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
)

type (
	Options struct {
		Logger logger.Logger
		DB     db.ExtendedDB
	}
	Container interface {
		Account() account.Repository
		Transfer() transfer.Repository
	}
	container struct {
		account  account.Repository
		transfer transfer.Repository
	}
)

func NewContainer(opts Options) Container {
	return &container{
		account: account.NewRepository(account.Options{
			Logger: opts.Logger,
			DB:     opts.DB,
		}),
		transfer: transfer.NewRepository(transfer.Options{
			Logger: opts.Logger,
			DB:     opts.DB,
		}),
	}
}

func (c *container) Account() account.Repository {
	return c.account
}

func (c *container) Transfer() transfer.Repository {
	return c.transfer
}
