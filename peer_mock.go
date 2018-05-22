// Code generated by MockGen. DO NOT EDIT.
// Source: peer.go

// Package gohfc is a generated GoMock package.
package gohfc

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	peer "github.com/hyperledger/fabric/protos/peer"
	grpc "google.golang.org/grpc"
	reflect "reflect"
)

// MockPeer is a mock of Peer interface
type MockPeer struct {
	ctrl     *gomock.Controller
	recorder *MockPeerMockRecorder
}

// MockPeerMockRecorder is the mock recorder for MockPeer
type MockPeerMockRecorder struct {
	mock *MockPeer
}

// NewMockPeer creates a new mock instance
func NewMockPeer(ctrl *gomock.Controller) *MockPeer {
	mock := &MockPeer{ctrl: ctrl}
	mock.recorder = &MockPeerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPeer) EXPECT() *MockPeerMockRecorder {
	return m.recorder
}

// Endorse mocks base method
func (m *MockPeer) Endorse(ctx context.Context, resp chan *PeerResponse, prop *peer.SignedProposal) {
	m.ctrl.Call(m, "Endorse", ctx, resp, prop)
}

// Endorse indicates an expected call of Endorse
func (mr *MockPeerMockRecorder) Endorse(ctx, resp, prop interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Endorse", reflect.TypeOf((*MockPeer)(nil).Endorse), ctx, resp, prop)
}

// GetURI mocks base method
func (m *MockPeer) GetURI() string {
	ret := m.ctrl.Call(m, "GetURI")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetURI indicates an expected call of GetURI
func (mr *MockPeerMockRecorder) GetURI() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetURI", reflect.TypeOf((*MockPeer)(nil).GetURI))
}

// GetOpts mocks base method
func (m *MockPeer) GetOpts() []grpc.DialOption {
	ret := m.ctrl.Call(m, "GetOpts")
	ret0, _ := ret[0].([]grpc.DialOption)
	return ret0
}

// GetOpts indicates an expected call of GetOpts
func (mr *MockPeerMockRecorder) GetOpts() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpts", reflect.TypeOf((*MockPeer)(nil).GetOpts))
}
