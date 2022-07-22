// Code generated by MockGen. DO NOT EDIT.
// Source: transfer.go

// Package transfer is a generated GoMock package.
package transfer

import (
	context "context"
	reflect "reflect"

	transaction "github.com/carlosrodriguesf/bank-api/pkg/apputil/transaction"
	model "github.com/carlosrodriguesf/bank-api/pkg/model"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(ctx context.Context, movement model.Transfer) (*model.GeneratedData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, movement)
	ret0, _ := ret[0].(*model.GeneratedData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(ctx, movement interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), ctx, movement)
}

// CreateMovement mocks base method.
func (m *MockRepository) CreateMovement(ctx context.Context, transferID, movementID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMovement", ctx, transferID, movementID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateMovement indicates an expected call of CreateMovement.
func (mr *MockRepositoryMockRecorder) CreateMovement(ctx, transferID, movementID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMovement", reflect.TypeOf((*MockRepository)(nil).CreateMovement), ctx, transferID, movementID)
}

// WithTransaction mocks base method.
func (m *MockRepository) WithTransaction(conn transaction.Transaction) Repository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", conn)
	ret0, _ := ret[0].(Repository)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockRepositoryMockRecorder) WithTransaction(conn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockRepository)(nil).WithTransaction), conn)
}