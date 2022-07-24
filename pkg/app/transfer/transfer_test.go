package transfer

import (
	"context"
	"errors"
	pkgerror "github.com/carlosrodriguesf/bank-api/pkg/error"
	"github.com/carlosrodriguesf/bank-api/pkg/model"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/account"
	"github.com/carlosrodriguesf/bank-api/pkg/repository/transfer"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/logger"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/transaction"
	"github.com/carlosrodriguesf/bank-api/pkg/tool/validator"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestList(t *testing.T) {
	transfersExample := []model.TransferDetailed{{
		Transfer: model.Transfer{
			ID:              "transfer_id",
			OriginAccountID: "origin_account_id",
			TargetAccountID: "target_account_id",
			Amount:          500,
			CreatedAt:       time.Now(),
		},
		OriginAccountName: "Origin Account",
		TargetAccountName: "Target account",
	}}
	cases := map[string]struct {
		InputData               string
		ExpectedData            []model.TransferDetailed
		ExpectedError           error
		PrepareMockRepoTransfer func(mock *transfer.MockRepository)
	}{
		"should return success": {
			InputData:     "origin_account_id",
			ExpectedData:  transfersExample,
			ExpectedError: nil,
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository) {
				mock.EXPECT().List(gomock.Any(), "origin_account_id").Return(transfersExample, nil)
			},
		},
		"should return error": {
			InputData:     "origin_account_id",
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantListTransfers,
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository) {
				mock.EXPECT().List(gomock.Any(), "origin_account_id").Return(nil, errors.New("fail"))
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl, ctx        = gomock.WithContext(context.Background(), t)
				mockRepoTransfer = transfer.NewMockRepository(ctrl)
				app              = NewApp(Options{
					Logger:       logger.New(""),
					RepoTransfer: mockRepoTransfer,
				})
			)

			cs.PrepareMockRepoTransfer(mockRepoTransfer)

			data, err := app.List(ctx, cs.InputData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}

func TestCreate(t *testing.T) {
	var (
		currentTime = time.Now()
		createData  = model.Transfer{
			OriginAccountID: "origin_account_id",
			TargetAccountID: "target_account_id",
			Amount:          500,
		}
		genTransferData = model.GeneratedData{
			ID:        "transfer_id",
			CreatedAt: currentTime,
		}
		accountOrigin = model.Account{
			ID:       createData.OriginAccountID,
			Name:     "Test Account",
			Document: "12312312312",
			Balance:  1000,
		}
		accountTarget = model.Account{
			ID:       createData.TargetAccountID,
			Name:     "Test Account",
			Document: "12312312312",
			Balance:  1000,
		}
		createdTransfer = model.Transfer{
			ID:              genTransferData.ID,
			OriginAccountID: createData.OriginAccountID,
			TargetAccountID: createData.TargetAccountID,
			Amount:          createData.Amount,
			CreatedAt:       genTransferData.CreatedAt,
		}
		validationError = validator.ValidationError{}
	)
	cases := map[string]struct {
		InputData               model.Transfer
		ExpectedData            *model.Transfer
		ExpectedError           error
		PrepareMockValidator    func(mock *validator.MockValidator)
		PrepareMockTxManager    func(mock *transaction.MockManager, tx transaction.Transaction)
		PrepareMockRepoAccount  func(mock *account.MockRepository, tx transaction.Transaction)
		PrepareMockRepoTransfer func(mock *transfer.MockRepository, tx transaction.Transaction)
	}{
		"should return success": {
			InputData:     createData,
			ExpectedData:  &createdTransfer,
			ExpectedError: nil,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(tx, nil)
				mock.EXPECT().Commit(tx)
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountOrigin.ID,
						accountOrigin.Balance-createData.Amount).
					Return(nil)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountTarget.ID,
						accountOrigin.Balance+createData.Amount).
					Return(nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().Create(gomock.Any(), createData).Return(&genTransferData, nil)
			},
		},
		"should return error: validation": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: &validationError,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(&validationError)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: can't get origin account": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(nil, errors.New("fail"))
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: origin account not exists": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrOriginAccountTransferNotFound,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(nil, nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: can't get target account": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(nil, errors.New("fail"))
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: origin target not exists": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrTargetAccountTransferNotFound,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(nil, nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: insufficient funds": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrInsufficientFunds,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				accountOrigin := accountOrigin
				accountOrigin.Balance = 0
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: create transaction": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(nil, errors.New("fail"))
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
			},
		},
		"should return error: can't create transfer": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(tx, nil)
				mock.EXPECT().Rollback(tx)
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountOrigin.ID,
						accountOrigin.Balance-createData.Amount).
					Return(nil)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountTarget.ID,
						accountOrigin.Balance+createData.Amount).
					Return(nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().Create(gomock.Any(), createData).Return(nil, errors.New("fail"))
			},
		},
		"should return error: can't update origin account balance": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(tx, nil)
				mock.EXPECT().Rollback(tx)
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountOrigin.ID,
						accountOrigin.Balance-createData.Amount).
					Return(errors.New("fail"))
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountTarget.ID,
						accountOrigin.Balance+createData.Amount).
					Return(nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().Create(gomock.Any(), createData).Return(&genTransferData, nil)
			},
		},
		"should return error: can't update target account balance": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(tx, nil)
				mock.EXPECT().Rollback(tx)
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountOrigin.ID,
						accountOrigin.Balance-createData.Amount).
					Return(nil)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountTarget.ID,
						accountOrigin.Balance+createData.Amount).
					Return(errors.New("fail"))
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().Create(gomock.Any(), createData).Return(&genTransferData, nil)
			},
		},
		"should return error: can't commit transaction": {
			InputData:     createData,
			ExpectedData:  nil,
			ExpectedError: pkgerror.ErrCantCreateTransfer,
			PrepareMockValidator: func(mock *validator.MockValidator) {
				mock.EXPECT().Validate(createData).Return(nil)
			},
			PrepareMockTxManager: func(mock *transaction.MockManager, tx transaction.Transaction) {
				mock.EXPECT().Create(gomock.Any()).Return(tx, nil)
				mock.EXPECT().Commit(tx).Return(errors.New("fail"))
				mock.EXPECT().Rollback(tx).Return(errors.New("fail"))
			},
			PrepareMockRepoAccount: func(mock *account.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.OriginAccountID).Return(&accountOrigin, nil)
				mock.EXPECT().GetByIDOrDocument(gomock.Any(), createData.TargetAccountID).Return(&accountTarget, nil)
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountOrigin.ID,
						accountOrigin.Balance-createData.Amount).
					Return(nil)
				mock.EXPECT().
					UpdateBalance(
						gomock.Any(),
						accountTarget.ID,
						accountOrigin.Balance+createData.Amount).
					Return(nil)
			},
			PrepareMockRepoTransfer: func(mock *transfer.MockRepository, tx transaction.Transaction) {
				mock.EXPECT().WithTransaction(tx).Return(mock)
				mock.EXPECT().Create(gomock.Any(), createData).Return(&genTransferData, nil)
			},
		},
	}

	for name, cs := range cases {
		t.Run(name, func(t *testing.T) {
			var (
				ctrl, ctx        = gomock.WithContext(context.Background(), t)
				txExample        = transaction.Transaction(nil)
				mockValidator    = validator.NewMockValidator(ctrl)
				mockTxManager    = transaction.NewMockManager(ctrl)
				mockRepoAccount  = account.NewMockRepository(ctrl)
				mockRepoTransfer = transfer.NewMockRepository(ctrl)
				app              = NewApp(Options{
					Logger:       logger.New(""),
					Validator:    mockValidator,
					TxManager:    mockTxManager,
					RepoAccount:  mockRepoAccount,
					RepoTransfer: mockRepoTransfer,
				})
			)

			cs.PrepareMockValidator(mockValidator)
			cs.PrepareMockTxManager(mockTxManager, txExample)
			cs.PrepareMockRepoAccount(mockRepoAccount, txExample)
			cs.PrepareMockRepoTransfer(mockRepoTransfer, txExample)

			data, err := app.Create(ctx, createData)

			assert.Equal(t, cs.ExpectedData, data)
			assert.Equal(t, cs.ExpectedError, err)
		})
	}
}
