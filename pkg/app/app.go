package app

import (
	"github.com/carlosrodriguesf/bank-api/pkg/app/account"
	"github.com/carlosrodriguesf/bank-api/pkg/repository"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
)

type (
	Options struct {
		Repository repository.Container
		Logger     logger.Logger
	}
	Container interface {
		Account() account.App
	}
	container struct {
		account account.App
	}
)

func NewContainer(opts Options) Container {
	var (
		validatorInstance = validator.New()
		secretInstance    = secret.New()
	)
	return &container{
		account: account.NewApp(account.Options{
			RepoAccount: opts.Repository.Account(),
			Logger:      opts.Logger,
			Validator:   validatorInstance,
			Secret:      secretInstance,
		}),
	}
}

func (c *container) Account() account.App {
	return c.account
}
