// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/RomanSaveljev/android-symbols/transmitter/receiver (interfaces: Receiver)

package mock

import (
	signatures "github.com/RomanSaveljev/android-symbols/transmitter/signatures"
	gomock "github.com/golang/mock/gomock"
)

// Mock of Receiver interface
type MockReceiver struct {
	ctrl     *gomock.Controller
	recorder *_MockReceiverRecorder
}

// Recorder for MockReceiver (not exported)
type _MockReceiverRecorder struct {
	mock *MockReceiver
}

func NewMockReceiver(ctrl *gomock.Controller) *MockReceiver {
	mock := &MockReceiver{ctrl: ctrl}
	mock.recorder = &_MockReceiverRecorder{mock}
	return mock
}

func (_m *MockReceiver) EXPECT() *_MockReceiverRecorder {
	return _m.recorder
}

func (_m *MockReceiver) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockReceiverRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockReceiver) SaveChunk(_param0 uint32, _param1 []byte, _param2 []byte) error {
	ret := _m.ctrl.Call(_m, "SaveChunk", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockReceiverRecorder) SaveChunk(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SaveChunk", arg0, arg1, arg2)
}

func (_m *MockReceiver) Signatures() (signatures.Signatures, error) {
	ret := _m.ctrl.Call(_m, "Signatures")
	ret0, _ := ret[0].(signatures.Signatures)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockReceiverRecorder) Signatures() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Signatures")
}

func (_m *MockReceiver) Write(_param0 []byte) (int, error) {
	ret := _m.ctrl.Call(_m, "Write", _param0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockReceiverRecorder) Write(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Write", arg0)
}
