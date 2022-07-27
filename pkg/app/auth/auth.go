//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package auth

import (
	"context"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/cache"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/generate"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"github.com/google/uuid"
	"time"
)

const (
	cacheKeySession = "auth:session:%s"
	cacheExpiration = time.Hour
)

type (
	Options struct {
		Logger      logger.Logger
		Secret      secret.Secret
		Cache       cache.Cache
		Validator   validator.Validator
		RepoAccount account.Repository
		Generate    generate.Generate
	}
	App interface {
		Auth(ctx context.Context, credentials model.Credentials) (*model.Session, error)
		GetSessionByToken(ctx context.Context, token string) (*model.Session, error)
	}
	appImpl struct {
		logger      logger.Logger
		secret      secret.Secret
		cache       cache.Cache
		validator   validator.Validator
		repoAccount account.Repository
		generate    generate.Generate
	}
)

func NewApp(opts Options) App {
	return &appImpl{
		logger:      opts.Logger.WithLocation().WithPreffix("service.auth"),
		secret:      opts.Secret,
		cache:       opts.Cache,
		validator:   opts.Validator,
		repoAccount: opts.RepoAccount,
	}
}

func (a *appImpl) Auth(ctx context.Context, credentials model.Credentials) (*model.Session, error) {
	if err := a.validator.Validate(credentials); err != nil {
		return nil, err
	}

	credentials.Document = model.DocumentRegex.ReplaceAllString(credentials.Document, "")

	acc, err := a.repoAccount.GetByIDOrDocument(ctx, credentials.Document)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantAuth
	}
	if acc == nil {
		return nil, pkgerror.ErrInvalidCredentials
	}
	if !a.secret.Verify(credentials.Secret, acc.Secret, acc.SecretSalt) {
		return nil, pkgerror.ErrInvalidCredentials
	}

	session := &model.Session{
		Token:     uuid.NewString(),
		Account:   *acc,
		CreatedAt: time.Now(),
	}

	err = a.cache.Set(ctx, getSessionCacheKey(session.Token), session, cacheExpiration)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantAuth
	}

	return session, nil
}

func (a *appImpl) GetSessionByToken(ctx context.Context, token string) (*model.Session, error) {
	var (
		session = new(model.Session)
		err     = a.cache.GetUpdating(ctx, getSessionCacheKey(token), session, cacheExpiration)
	)
	if err != nil {
		if a.cache.IsErrCacheMissing(err) {
			return nil, pkgerror.ErrSessionNotFound
		}
		a.logger.Error(err)
		return nil, pkgerror.ErrCantGetSession
	}
	return session, nil
}
