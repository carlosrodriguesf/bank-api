//go:generate mockgen -source=${GOFILE} -package=${GOPACKAGE} -destination=${GOPACKAGE}_mock.go

package transfer

import (
	"context"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/transfer"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/transaction"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"golang.org/x/sync/errgroup"
)

type (
	Options struct {
		Logger       logger.Logger
		Validator    validator.Validator
		TxManager    transaction.Manager
		RepoAccount  account.Repository
		RepoTransfer transfer.Repository
	}
	App interface {
		Create(ctx context.Context, transfer model.Transfer) (*model.Transfer, error)
		List(ctx context.Context, accountID string) ([]model.TransferDetailed, error)
	}
	appImpl struct {
		logger       logger.Logger
		validator    validator.Validator
		txManager    transaction.Manager
		repoAccount  account.Repository
		repoTransfer transfer.Repository
	}
)

func NewApp(opts Options) App {
	return &appImpl{
		logger:       opts.Logger.WithLocation().WithPreffix("app.transfer"),
		validator:    opts.Validator,
		txManager:    opts.TxManager,
		repoAccount:  opts.RepoAccount,
		repoTransfer: opts.RepoTransfer,
	}
}

func (a *appImpl) List(ctx context.Context, accountID string) ([]model.TransferDetailed, error) {
	transfers, err := a.repoTransfer.List(ctx, accountID)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantListTransfers
	}
	return transfers, nil
}

func (a appImpl) Create(ctx context.Context, transfer model.Transfer) (*model.Transfer, error) {
	if err := a.validator.Validate(transfer); err != nil {
		return nil, err
	}

	originAccount, err := a.repoAccount.GetByIDOrDocument(ctx, transfer.OriginAccountID)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantCreateTransfer
	}
	if originAccount == nil {
		return nil, pkgerror.ErrOriginAccountTransferNotFound
	}
	if originAccount.Balance < transfer.Amount {
		return nil, pkgerror.ErrInsufficientFunds
	}

	targetAccount, err := a.repoAccount.GetByIDOrDocument(ctx, transfer.TargetAccountID)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantCreateTransfer
	}
	if targetAccount == nil {
		return nil, pkgerror.ErrTargetAccountTransferNotFound
	}

	tx, err := a.startTransaction(ctx)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantCreateTransfer
	}
	defer func() {
		if err != nil {
			a.rollbackTransaction(tx)
		}
	}()

	genData, err := a.makeTransfer(ctx, transferWrapper{
		Transfer:      transfer,
		AccountOrigin: originAccount,
		AccountTarget: targetAccount,
	})
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantCreateTransfer
	}

	transfer.ID = genData.ID
	transfer.CreatedAt = genData.CreatedAt

	err = a.txManager.Commit(tx)
	if err != nil {
		a.logger.Error(err)
		return nil, pkgerror.ErrCantCreateTransfer
	}

	return &transfer, nil
}

func (a *appImpl) makeTransfer(ctx context.Context, wrapper transferWrapper) (*model.GeneratedData, error) {
	var (
		transferData  = wrapper.Transfer
		accountOrigin = wrapper.AccountOrigin
		accountTarget = wrapper.AccountTarget
		errGroup      *errgroup.Group
		genData       *model.GeneratedData
	)

	errGroup, ctx = errgroup.WithContext(ctx)

	errGroup.Go(func() (err error) {
		genData, err = a.repoTransfer.Create(ctx, transferData)
		if err != nil {
			a.logger.Error(err)
		}
		return
	})

	errGroup.Go(func() (err error) {
		err = a.repoAccount.UpdateBalance(ctx, accountOrigin.ID, accountOrigin.Balance-transferData.Amount)
		if err != nil {
			a.logger.Error(err)
		}
		return
	})

	errGroup.Go(func() (err error) {
		err = a.repoAccount.UpdateBalance(ctx, accountTarget.ID, accountTarget.Balance+transferData.Amount)
		if err != nil {
			a.logger.Error(err)
		}
		return err
	})

	if err := errGroup.Wait(); err != nil {
		a.logger.Error(err)
		return nil, err
	}

	return genData, nil
}

func (a *appImpl) startTransaction(ctx context.Context) (transaction.Transaction, error) {
	tx, err := a.txManager.Create(ctx)
	if err != nil {
		a.logger.Error(err)
		return nil, err
	}

	a.repoAccount = a.repoAccount.WithTransaction(tx)
	a.repoTransfer = a.repoTransfer.WithTransaction(tx)

	return tx, nil
}

func (a *appImpl) rollbackTransaction(tx transaction.Transaction) {
	err := a.txManager.Rollback(tx)
	if err != nil {
		a.logger.Error(err)
	}
}
