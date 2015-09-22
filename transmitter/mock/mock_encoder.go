// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/RomanSaveljev/android-symbols/transmitter/encoder (interfaces: Encoder)

package mock

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Encoder interface
type MockEncoder struct {
	ctrl     *gomock.Controller
	recorder *_MockEncoderRecorder
}

// Recorder for MockEncoder (not exported)
type _MockEncoderRecorder struct {
	mock *MockEncoder
}

func NewMockEncoder(ctrl *gomock.Controller) *MockEncoder {
	mock := &MockEncoder{ctrl: ctrl}
	mock.recorder = &_MockEncoderRecorder{mock}
	return mock
}

func (_m *MockEncoder) EXPECT() *_MockEncoderRecorder {
	return _m.recorder
}

func (_m *MockEncoder) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEncoderRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

func (_m *MockEncoder) Write(_param0 []byte) (int, error) {
	ret := _m.ctrl.Call(_m, "Write", _param0)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockEncoderRecorder) Write(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Write", arg0)
}

func (_m *MockEncoder) WriteSignature(_param0 uint32, _param1 []byte) error {
	ret := _m.ctrl.Call(_m, "WriteSignature", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockEncoderRecorder) WriteSignature(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "WriteSignature", arg0, arg1)
}
