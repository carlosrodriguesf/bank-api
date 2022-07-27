package app

import (
	"github.com/carlosrodriguesf/bank-api/pkg/app/account"
	"github.com/carlosrodriguesf/bank-api/pkg/app/auth"
	"github.com/carlosrodriguesf/bank-api/pkg/app/transfer"
	"github.com/carlosrodriguesf/bank-api/pkg/repository"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/cache"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/db"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/generate"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/transaction"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
)

type (
	Options struct {
		DB         db.ExtendedDB
		Repository repository.Container
		Logger     logger.Logger
		Cache      cache.Cache
	}
	Container interface {
		Account() account.App
		Auth() auth.App
		Transfer() transfer.App
	}
	container struct {
		account  account.App
		auth     auth.App
		transfer transfer.App
	}
)

func NewContainer(opts Options) Container {
	var (
		validatorInstance = validator.New()
		secretInstance    = secret.New()
		txManagerInstance = transaction.NewManager(opts.DB)
		generateInstance  = generate.New()
	)
	return &container{
		account: account.NewApp(account.Options{
			RepoAccount: opts.Repository.Account(),
			Logger:      opts.Logger,
			Validator:   validatorInstance,
			Secret:      secretInstance,
		}),
		auth: auth.NewApp(auth.Options{
			Logger:      opts.Logger,
			Cache:       opts.Cache,
			Validator:   validatorInstance,
			Secret:      secretInstance,
			RepoAccount: opts.Repository.Account(),
			Generate:    generateInstance,
		}),
		transfer: transfer.NewApp(transfer.Options{
			Logger:       opts.Logger,
			Validator:    validatorInstance,
			TxManager:    txManagerInstance,
			RepoAccount:  opts.Repository.Account(),
			RepoTransfer: opts.Repository.Transfer(),
		}),
	}
}

func (c *container) Account() account.App {
	return c.account
}

func (c *container) Auth() auth.App {
	return c.auth
}

func (c *container) Transfer() transfer.App {
	return c.transfer
}
