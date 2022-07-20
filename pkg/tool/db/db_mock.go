// Code generated by MockGen. DO NOT EDIT.
// Source: db.go

// Package db is a generated GoMock package.
package db

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
)

// MockExtendedDB is a mock of ExtendedDB interface.
type MockExtendedDB struct {
	ctrl     *gomock.Controller
	recorder *MockExtendedDBMockRecorder
}

// MockExtendedDBMockRecorder is the mock recorder for MockExtendedDB.
type MockExtendedDBMockRecorder struct {
	mock *MockExtendedDB
}

// NewMockExtendedDB creates a new mock instance.
func NewMockExtendedDB(ctrl *gomock.Controller) *MockExtendedDB {
	mock := &MockExtendedDB{ctrl: ctrl}
	mock.recorder = &MockExtendedDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExtendedDB) EXPECT() *MockExtendedDBMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockExtendedDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockExtendedDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockExtendedDB)(nil).Close))
}

// ExecContext mocks base method.
func (m *MockExtendedDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockExtendedDBMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockExtendedDB)(nil).ExecContext), varargs...)
}

// GetContext mocks base method.
func (m *MockExtendedDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockExtendedDBMockRecorder) GetContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockExtendedDB)(nil).GetContext), varargs...)
}

// NamedExecContext mocks base method.
func (m *MockExtendedDB) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedExecContext", ctx, query, arg)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamedExecContext indicates an expected call of NamedExecContext.
func (mr *MockExtendedDBMockRecorder) NamedExecContext(ctx, query, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedExecContext", reflect.TypeOf((*MockExtendedDB)(nil).NamedExecContext), ctx, query, arg)
}

// NamedGetContext mocks base method.
func (m *MockExtendedDB) NamedGetContext(ctx context.Context, query string, dest, arg interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedGetContext", ctx, query, dest, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// NamedGetContext indicates an expected call of NamedGetContext.
func (mr *MockExtendedDBMockRecorder) NamedGetContext(ctx, query, dest, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedGetContext", reflect.TypeOf((*MockExtendedDB)(nil).NamedGetContext), ctx, query, dest, arg)
}

// NamedQueryContext mocks base method.
func (m *MockExtendedDB) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedQueryContext", ctx, query, arg)
	ret0, _ := ret[0].(*sqlx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamedQueryContext indicates an expected call of NamedQueryContext.
func (mr *MockExtendedDBMockRecorder) NamedQueryContext(ctx, query, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedQueryContext", reflect.TypeOf((*MockExtendedDB)(nil).NamedQueryContext), ctx, query, arg)
}

// NamedSelectContext mocks base method.
func (m *MockExtendedDB) NamedSelectContext(ctx context.Context, query string, dest, arg interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedSelectContext", ctx, query, dest, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// NamedSelectContext indicates an expected call of NamedSelectContext.
func (mr *MockExtendedDBMockRecorder) NamedSelectContext(ctx, query, dest, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedSelectContext", reflect.TypeOf((*MockExtendedDB)(nil).NamedSelectContext), ctx, query, dest, arg)
}

// PrepareNamedContext mocks base method.
func (m *MockExtendedDB) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareNamedContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.NamedStmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareNamedContext indicates an expected call of PrepareNamedContext.
func (mr *MockExtendedDBMockRecorder) PrepareNamedContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareNamedContext", reflect.TypeOf((*MockExtendedDB)(nil).PrepareNamedContext), ctx, query)
}

// PreparexContext mocks base method.
func (m *MockExtendedDB) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreparexContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreparexContext indicates an expected call of PreparexContext.
func (mr *MockExtendedDBMockRecorder) PreparexContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreparexContext", reflect.TypeOf((*MockExtendedDB)(nil).PreparexContext), ctx, query)
}

// QueryRowxContext mocks base method.
func (m *MockExtendedDB) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowxContext", varargs...)
	ret0, _ := ret[0].(*sqlx.Row)
	return ret0
}

// QueryRowxContext indicates an expected call of QueryRowxContext.
func (mr *MockExtendedDBMockRecorder) QueryRowxContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowxContext", reflect.TypeOf((*MockExtendedDB)(nil).QueryRowxContext), varargs...)
}

// QueryxContext mocks base method.
func (m *MockExtendedDB) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryxContext", varargs...)
	ret0, _ := ret[0].(*sqlx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryxContext indicates an expected call of QueryxContext.
func (mr *MockExtendedDBMockRecorder) QueryxContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryxContext", reflect.TypeOf((*MockExtendedDB)(nil).QueryxContext), varargs...)
}

// SelectContext mocks base method.
func (m *MockExtendedDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectContext indicates an expected call of SelectContext.
func (mr *MockExtendedDBMockRecorder) SelectContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContext", reflect.TypeOf((*MockExtendedDB)(nil).SelectContext), varargs...)
}

// MockExtendedTx is a mock of ExtendedTx interface.
type MockExtendedTx struct {
	ctrl     *gomock.Controller
	recorder *MockExtendedTxMockRecorder
}

// MockExtendedTxMockRecorder is the mock recorder for MockExtendedTx.
type MockExtendedTxMockRecorder struct {
	mock *MockExtendedTx
}

// NewMockExtendedTx creates a new mock instance.
func NewMockExtendedTx(ctrl *gomock.Controller) *MockExtendedTx {
	mock := &MockExtendedTx{ctrl: ctrl}
	mock.recorder = &MockExtendedTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockExtendedTx) EXPECT() *MockExtendedTxMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockExtendedTx) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockExtendedTxMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockExtendedTx)(nil).Close))
}

// Commit mocks base method.
func (m *MockExtendedTx) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockExtendedTxMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockExtendedTx)(nil).Commit))
}

// ExecContext mocks base method.
func (m *MockExtendedTx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockExtendedTxMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockExtendedTx)(nil).ExecContext), varargs...)
}

// GetContext mocks base method.
func (m *MockExtendedTx) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockExtendedTxMockRecorder) GetContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockExtendedTx)(nil).GetContext), varargs...)
}

// NamedExecContext mocks base method.
func (m *MockExtendedTx) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedExecContext", ctx, query, arg)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamedExecContext indicates an expected call of NamedExecContext.
func (mr *MockExtendedTxMockRecorder) NamedExecContext(ctx, query, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedExecContext", reflect.TypeOf((*MockExtendedTx)(nil).NamedExecContext), ctx, query, arg)
}

// NamedGetContext mocks base method.
func (m *MockExtendedTx) NamedGetContext(ctx context.Context, query string, dest, arg interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedGetContext", ctx, query, dest, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// NamedGetContext indicates an expected call of NamedGetContext.
func (mr *MockExtendedTxMockRecorder) NamedGetContext(ctx, query, dest, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedGetContext", reflect.TypeOf((*MockExtendedTx)(nil).NamedGetContext), ctx, query, dest, arg)
}

// NamedQueryContext mocks base method.
func (m *MockExtendedTx) NamedQueryContext(ctx context.Context, query string, arg interface{}) (*sqlx.Rows, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedQueryContext", ctx, query, arg)
	ret0, _ := ret[0].(*sqlx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NamedQueryContext indicates an expected call of NamedQueryContext.
func (mr *MockExtendedTxMockRecorder) NamedQueryContext(ctx, query, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedQueryContext", reflect.TypeOf((*MockExtendedTx)(nil).NamedQueryContext), ctx, query, arg)
}

// NamedSelectContext mocks base method.
func (m *MockExtendedTx) NamedSelectContext(ctx context.Context, query string, dest, arg interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NamedSelectContext", ctx, query, dest, arg)
	ret0, _ := ret[0].(error)
	return ret0
}

// NamedSelectContext indicates an expected call of NamedSelectContext.
func (mr *MockExtendedTxMockRecorder) NamedSelectContext(ctx, query, dest, arg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NamedSelectContext", reflect.TypeOf((*MockExtendedTx)(nil).NamedSelectContext), ctx, query, dest, arg)
}

// PrepareNamedContext mocks base method.
func (m *MockExtendedTx) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PrepareNamedContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.NamedStmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PrepareNamedContext indicates an expected call of PrepareNamedContext.
func (mr *MockExtendedTxMockRecorder) PrepareNamedContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PrepareNamedContext", reflect.TypeOf((*MockExtendedTx)(nil).PrepareNamedContext), ctx, query)
}

// PreparexContext mocks base method.
func (m *MockExtendedTx) PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PreparexContext", ctx, query)
	ret0, _ := ret[0].(*sqlx.Stmt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PreparexContext indicates an expected call of PreparexContext.
func (mr *MockExtendedTxMockRecorder) PreparexContext(ctx, query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PreparexContext", reflect.TypeOf((*MockExtendedTx)(nil).PreparexContext), ctx, query)
}

// QueryRowxContext mocks base method.
func (m *MockExtendedTx) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryRowxContext", varargs...)
	ret0, _ := ret[0].(*sqlx.Row)
	return ret0
}

// QueryRowxContext indicates an expected call of QueryRowxContext.
func (mr *MockExtendedTxMockRecorder) QueryRowxContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryRowxContext", reflect.TypeOf((*MockExtendedTx)(nil).QueryRowxContext), varargs...)
}

// QueryxContext mocks base method.
func (m *MockExtendedTx) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryxContext", varargs...)
	ret0, _ := ret[0].(*sqlx.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryxContext indicates an expected call of QueryxContext.
func (mr *MockExtendedTxMockRecorder) QueryxContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryxContext", reflect.TypeOf((*MockExtendedTx)(nil).QueryxContext), varargs...)
}

// Rollback mocks base method.
func (m *MockExtendedTx) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockExtendedTxMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockExtendedTx)(nil).Rollback))
}

// SelectContext mocks base method.
func (m *MockExtendedTx) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectContext indicates an expected call of SelectContext.
func (mr *MockExtendedTxMockRecorder) SelectContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContext", reflect.TypeOf((*MockExtendedTx)(nil).SelectContext), varargs...)
}
