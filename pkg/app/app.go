package app

import (
	"github.com/carlosrodriguesf/bank-api/pkg/app/account"
	"github.com/carlosrodriguesf/bank-api/pkg/app/auth"
	"github.com/carlosrodriguesf/bank-api/pkg/repository"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/cache"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
)

type (
	Options struct {
		Repository repository.Container
		Logger     logger.Logger
		Cache      cache.Cache
	}
	Container interface {
		Account() account.App
		Auth() auth.App
	}
	container struct {
		account account.App
		auth    auth.App
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
		auth: auth.NewApp(auth.Options{
			Logger:      opts.Logger,
			Cache:       opts.Cache,
			Validator:   validatorInstance,
			Secret:      secretInstance,
			RepoAccount: opts.Repository.Account(),
		}),
	}
}

func (c *container) Account() account.App {
	return c.account
}

func (c *container) Auth() auth.App {
	return c.auth
}
