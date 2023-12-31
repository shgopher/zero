// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/superproj/zero/internal/gateway/biz (interfaces: BizFactory)

// Package biz is a generated GoMock package.
package biz

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	miner "github.com/superproj/zero/internal/gateway/biz/miner"
	minerset "github.com/superproj/zero/internal/gateway/biz/minerset"
)

// MockBizFactory is a mock of BizFactory interface.
type MockBizFactory struct {
	ctrl     *gomock.Controller
	recorder *MockBizFactoryMockRecorder
}

// MockBizFactoryMockRecorder is the mock recorder for MockBizFactory.
type MockBizFactoryMockRecorder struct {
	mock *MockBizFactory
}

// NewMockBizFactory creates a new mock instance.
func NewMockBizFactory(ctrl *gomock.Controller) *MockBizFactory {
	mock := &MockBizFactory{ctrl: ctrl}
	mock.recorder = &MockBizFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBizFactory) EXPECT() *MockBizFactoryMockRecorder {
	return m.recorder
}

// MinerSets mocks base method.
func (m *MockBizFactory) MinerSets() minerset.MinerSetBiz {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MinerSets")
	ret0, _ := ret[0].(minerset.MinerSetBiz)
	return ret0
}

// MinerSets indicates an expected call of MinerSets.
func (mr *MockBizFactoryMockRecorder) MinerSets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MinerSets", reflect.TypeOf((*MockBizFactory)(nil).MinerSets))
}

// Miners mocks base method.
func (m *MockBizFactory) Miners() miner.MinerBiz {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Miners")
	ret0, _ := ret[0].(miner.MinerBiz)
	return ret0
}

// Miners indicates an expected call of Miners.
func (mr *MockBizFactoryMockRecorder) Miners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Miners", reflect.TypeOf((*MockBizFactory)(nil).Miners))
}
