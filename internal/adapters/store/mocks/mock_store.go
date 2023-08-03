// Code generated by MockGen. DO NOT EDIT.
// Source: internal/adapters/store/store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/rawen554/go-loyal/internal/models"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStore) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// CreateUser mocks base method.
func (m *MockStore) CreateUser(user *models.User) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoreMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStore)(nil).CreateUser), user)
}

// CreateWithdraw mocks base method.
func (m *MockStore) CreateWithdraw(userID uint64, w models.BalanceWithdrawShema) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateWithdraw", userID, w)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateWithdraw indicates an expected call of CreateWithdraw.
func (mr *MockStoreMockRecorder) CreateWithdraw(userID, w interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateWithdraw", reflect.TypeOf((*MockStore)(nil).CreateWithdraw), userID, w)
}

// GetUnprocessedOrders mocks base method.
func (m *MockStore) GetUnprocessedOrders() ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUnprocessedOrders")
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUnprocessedOrders indicates an expected call of GetUnprocessedOrders.
func (mr *MockStoreMockRecorder) GetUnprocessedOrders() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUnprocessedOrders", reflect.TypeOf((*MockStore)(nil).GetUnprocessedOrders))
}

// GetUser mocks base method.
func (m *MockStore) GetUser(u *models.User) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", u)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(u interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), u)
}

// GetUserBalance mocks base method.
func (m *MockStore) GetUserBalance(userID uint64) (*models.UserBalanceShema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", userID)
	ret0, _ := ret[0].(*models.UserBalanceShema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockStoreMockRecorder) GetUserBalance(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockStore)(nil).GetUserBalance), userID)
}

// GetUserOrders mocks base method.
func (m *MockStore) GetUserOrders(userID uint64) ([]models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserOrders", userID)
	ret0, _ := ret[0].([]models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserOrders indicates an expected call of GetUserOrders.
func (mr *MockStoreMockRecorder) GetUserOrders(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserOrders", reflect.TypeOf((*MockStore)(nil).GetUserOrders), userID)
}

// GetWithdrawals mocks base method.
func (m *MockStore) GetWithdrawals(userID uint64) ([]models.Withdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", userID)
	ret0, _ := ret[0].([]models.Withdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockStoreMockRecorder) GetWithdrawals(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockStore)(nil).GetWithdrawals), userID)
}

// Ping mocks base method.
func (m *MockStore) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoreMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStore)(nil).Ping))
}

// PutOrder mocks base method.
func (m *MockStore) PutOrder(number string, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutOrder", number, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutOrder indicates an expected call of PutOrder.
func (mr *MockStoreMockRecorder) PutOrder(number, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutOrder", reflect.TypeOf((*MockStore)(nil).PutOrder), number, userID)
}

// UpdateOrder mocks base method.
func (m *MockStore) UpdateOrder(o *models.Order) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", o)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStoreMockRecorder) UpdateOrder(o interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStore)(nil).UpdateOrder), o)
}