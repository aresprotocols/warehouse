// Code generated by MockGen. DO NOT EDIT.
// Source: price_api/price_server/internal/repository (interfaces: LogInfoRepository)

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	vo "price_api/price_server/internal/vo"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLogInfoRepository is a mock of LogInfoRepository interface.
type MockLogInfoRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLogInfoRepositoryMockRecorder
}

// MockLogInfoRepositoryMockRecorder is the mock recorder for MockLogInfoRepository.
type MockLogInfoRepositoryMockRecorder struct {
	mock *MockLogInfoRepository
}

// NewMockLogInfoRepository creates a new mock instance.
func NewMockLogInfoRepository(ctrl *gomock.Controller) *MockLogInfoRepository {
	mock := &MockLogInfoRepository{ctrl: ctrl}
	mock.recorder = &MockLogInfoRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLogInfoRepository) EXPECT() *MockLogInfoRepositoryMockRecorder {
	return m.recorder
}

// GetLogInfo mocks base method.
func (m *MockLogInfoRepository) GetLogInfo(arg0, arg1 int) (vo.LOG_INFOS, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogInfo", arg0, arg1)
	ret0, _ := ret[0].(vo.LOG_INFOS)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogInfo indicates an expected call of GetLogInfo.
func (mr *MockLogInfoRepositoryMockRecorder) GetLogInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogInfo", reflect.TypeOf((*MockLogInfoRepository)(nil).GetLogInfo), arg0, arg1)
}

// GetLogInfoBySymbol mocks base method.
func (m *MockLogInfoRepository) GetLogInfoBySymbol(arg0, arg1 int, arg2, arg3 string) ([]vo.REQ_RSP_LOG_INFO, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogInfoBySymbol", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]vo.REQ_RSP_LOG_INFO)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogInfoBySymbol indicates an expected call of GetLogInfoBySymbol.
func (mr *MockLogInfoRepositoryMockRecorder) GetLogInfoBySymbol(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogInfoBySymbol", reflect.TypeOf((*MockLogInfoRepository)(nil).GetLogInfoBySymbol), arg0, arg1, arg2, arg3)
}

// GetTotalLogInfoBySymbol mocks base method.
func (m *MockLogInfoRepository) GetTotalLogInfoBySymbol(arg0, arg1 string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotalLogInfoBySymbol", arg0, arg1)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTotalLogInfoBySymbol indicates an expected call of GetTotalLogInfoBySymbol.
func (mr *MockLogInfoRepositoryMockRecorder) GetTotalLogInfoBySymbol(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotalLogInfoBySymbol", reflect.TypeOf((*MockLogInfoRepository)(nil).GetTotalLogInfoBySymbol), arg0, arg1)
}

// InsertLogInfo mocks base method.
func (m *MockLogInfoRepository) InsertLogInfo(arg0 map[string]interface{}, arg1 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertLogInfo", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertLogInfo indicates an expected call of InsertLogInfo.
func (mr *MockLogInfoRepositoryMockRecorder) InsertLogInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertLogInfo", reflect.TypeOf((*MockLogInfoRepository)(nil).InsertLogInfo), arg0, arg1)
}
