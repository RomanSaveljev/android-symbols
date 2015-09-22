// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/RomanSaveljev/android-symbols/transmitter/chunker (interfaces: Chunker)

package mock

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Chunker interface
type MockChunker struct {
	ctrl     *gomock.Controller
	recorder *_MockChunkerRecorder
}

// Recorder for MockChunker (not exported)
type _MockChunkerRecorder struct {
	mock *MockChunker
}

func NewMockChunker(ctrl *gomock.Controller) *MockChunker {
	mock := &MockChunker{ctrl: ctrl}
	mock.recorder = &_MockChunkerRecorder{mock}
	return mock
}

func (_m *MockChunker) EXPECT() *_MockChunkerRecorder {
	return _m.recorder
}

func (_m *MockChunker) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockChunkerRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockChunker) Flush() error {
	ret := _m.ctrl.Call(_m, "Flush")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockChunkerRecorder) Flush() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Flush")
}

func (_m *MockChunker) Write(_param0 byte) error {
	ret := _m.ctrl.Call(_m, "Write", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockChunkerRecorder) Write(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Write", arg0)
}

func (_m *MockChunker) WriteSignature(_param0 uint32, _param1 []byte) error {
	ret := _m.ctrl.Call(_m, "WriteSignature", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockChunkerRecorder) WriteSignature(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WriteSignature", arg0, arg1)
}
