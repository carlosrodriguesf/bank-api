//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package account

import (
	"context"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/secret"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
)

type (
	Options struct {
		Logger      logger.Logger
		Secret      secret.Secret
		Validator   validator.Validator
		RepoAccount account.Repository
	}
	App interface {
		Create(ctx context.Context, account model.Account) (*model.Account, error)
	}
	appImpl struct {
		logger      logger.Logger
		secret      secret.Secret
		validator   validator.Validator
		repoAccount account.Repository
	}
)

func NewApp(opts Options) App {
	return &appImpl{
		logger:      opts.Logger.WithLocation().WithPreffix("service.account"),
		secret:      opts.Secret,
		validator:   opts.Validator,
		repoAccount: opts.RepoAccount,
	}
}

func (s *appImpl) Create(ctx context.Context, creationData model.Account) (*model.Account, error) {
	if err := s.validator.Validate(creationData); err != nil {
		return nil, err
	}

	creationData.SecretSalt = s.secret.GenSalt()
	creationData.Secret = s.secret.Encode(creationData.Secret, creationData.SecretSalt)
	creationData.Document = model.DocumentRegex.ReplaceAllString(creationData.Document, "")

	documentExists, err := s.repoAccount.HasDocument(ctx, creationData.Document)
	if err != nil {
		s.logger.Error(err)
		return nil, pkgerror.ErrCantCreateAccount
	}
	if documentExists {
		return nil, pkgerror.ErrDocumentAlreadyExists
	}

	generatedData, err := s.repoAccount.Create(ctx, creationData)
	if err != nil {
		s.logger.Error(err)
		return nil, pkgerror.ErrCantCreateAccount
	}

	return &model.Account{
		ID:         generatedData.ID,
		Name:       creationData.Name,
		Document:   creationData.Document,
		Secret:     creationData.Secret,
		SecretSalt: creationData.SecretSalt,
		CreatedAt:  generatedData.CreatedAt,
	}, nil
}
