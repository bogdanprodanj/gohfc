// Code generated by MockGen. DO NOT EDIT.
// Source: orderer.go

// Package gohfc is a generated GoMock package.
package gohfc

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	common "github.com/hyperledger/fabric/protos/common"
	orderer "github.com/hyperledger/fabric/protos/orderer"
	reflect "reflect"
)

// MockOrderer is a mock of Orderer interface
type MockOrderer struct {
	ctrl     *gomock.Controller
	recorder *MockOrdererMockRecorder
}

// MockOrdererMockRecorder is the mock recorder for MockOrderer
type MockOrdererMockRecorder struct {
	mock *MockOrderer
}

// NewMockOrderer creates a new mock instance
func NewMockOrderer(ctrl *gomock.Controller) *MockOrderer {
	mock := &MockOrderer{ctrl: ctrl}
	mock.recorder = &MockOrdererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOrderer) EXPECT() *MockOrdererMockRecorder {
	return m.recorder
}

// Broadcast mocks base method
func (m *MockOrderer) Broadcast(envelope *common.Envelope) (*orderer.BroadcastResponse, error) {
	ret := m.ctrl.Call(m, "Broadcast", envelope)
	ret0, _ := ret[0].(*orderer.BroadcastResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Broadcast indicates an expected call of Broadcast
func (mr *MockOrdererMockRecorder) Broadcast(envelope interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Broadcast", reflect.TypeOf((*MockOrderer)(nil).Broadcast), envelope)
}

// Deliver mocks base method
func (m *MockOrderer) Deliver(ctx context.Context, envelope *common.Envelope) (*common.Block, error) {
	ret := m.ctrl.Call(m, "Deliver", ctx, envelope)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Deliver indicates an expected call of Deliver
func (mr *MockOrdererMockRecorder) Deliver(ctx, envelope interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Deliver", reflect.TypeOf((*MockOrderer)(nil).Deliver), ctx, envelope)
}

// getGenesisBlock mocks base method
func (m *MockOrderer) getGenesisBlock(ctx context.Context, identity Identity, crypto CryptoSuite, channelId string) (*common.Block, error) {
	ret := m.ctrl.Call(m, "getGenesisBlock", ctx, identity, crypto, channelId)
	ret0, _ := ret[0].(*common.Block)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// getGenesisBlock indicates an expected call of getGenesisBlock
func (mr *MockOrdererMockRecorder) getGenesisBlock(ctx, identity, crypto, channelId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getGenesisBlock", reflect.TypeOf((*MockOrderer)(nil).getGenesisBlock), ctx, identity, crypto, channelId)
}
