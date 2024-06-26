// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/nickcorin/toolkit/flux (interfaces: CursorReader,CursorWriter,CursorStore)
//
// Generated by this command:
//
//	mockgen -write_generate_directive -write_package_comment -write_source_comment -package mocks -destination cursor.go github.com/nickcorin/toolkit/flux CursorReader,CursorWriter,CursorStore
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	flux "github.com/nickcorin/toolkit/flux"
	gomock "go.uber.org/mock/gomock"
)

//go:generate mockgen -write_generate_directive -write_package_comment -write_source_comment -package mocks -destination cursor.go github.com/nickcorin/toolkit/flux CursorReader,CursorWriter,CursorStore

// MockCursorReader is a mock of CursorReader interface.
type MockCursorReader struct {
	ctrl     *gomock.Controller
	recorder *MockCursorReaderMockRecorder
}

// MockCursorReaderMockRecorder is the mock recorder for MockCursorReader.
type MockCursorReaderMockRecorder struct {
	mock *MockCursorReader
}

// NewMockCursorReader creates a new mock instance.
func NewMockCursorReader(ctrl *gomock.Controller) *MockCursorReader {
	mock := &MockCursorReader{ctrl: ctrl}
	mock.recorder = &MockCursorReaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCursorReader) EXPECT() *MockCursorReaderMockRecorder {
	return m.recorder
}

// LookupCursorByID mocks base method.
func (m *MockCursorReader) LookupCursorByID(arg0 context.Context, arg1 string) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCursorByID", arg0, arg1)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCursorByID indicates an expected call of LookupCursorByID.
func (mr *MockCursorReaderMockRecorder) LookupCursorByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCursorByID", reflect.TypeOf((*MockCursorReader)(nil).LookupCursorByID), arg0, arg1)
}

// LookupCursorByName mocks base method.
func (m *MockCursorReader) LookupCursorByName(arg0 context.Context, arg1 string) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCursorByName", arg0, arg1)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCursorByName indicates an expected call of LookupCursorByName.
func (mr *MockCursorReaderMockRecorder) LookupCursorByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCursorByName", reflect.TypeOf((*MockCursorReader)(nil).LookupCursorByName), arg0, arg1)
}

// MockCursorWriter is a mock of CursorWriter interface.
type MockCursorWriter struct {
	ctrl     *gomock.Controller
	recorder *MockCursorWriterMockRecorder
}

// MockCursorWriterMockRecorder is the mock recorder for MockCursorWriter.
type MockCursorWriterMockRecorder struct {
	mock *MockCursorWriter
}

// NewMockCursorWriter creates a new mock instance.
func NewMockCursorWriter(ctrl *gomock.Controller) *MockCursorWriter {
	mock := &MockCursorWriter{ctrl: ctrl}
	mock.recorder = &MockCursorWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCursorWriter) EXPECT() *MockCursorWriterMockRecorder {
	return m.recorder
}

// CreateCursor mocks base method.
func (m *MockCursorWriter) CreateCursor(arg0 context.Context, arg1 string, arg2 uint) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCursor", arg0, arg1, arg2)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCursor indicates an expected call of CreateCursor.
func (mr *MockCursorWriterMockRecorder) CreateCursor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCursor", reflect.TypeOf((*MockCursorWriter)(nil).CreateCursor), arg0, arg1, arg2)
}

// UpdateCursor mocks base method.
func (m *MockCursorWriter) UpdateCursor(arg0 context.Context, arg1 string, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCursor", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCursor indicates an expected call of UpdateCursor.
func (mr *MockCursorWriterMockRecorder) UpdateCursor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCursor", reflect.TypeOf((*MockCursorWriter)(nil).UpdateCursor), arg0, arg1, arg2)
}

// MockCursorStore is a mock of CursorStore interface.
type MockCursorStore struct {
	ctrl     *gomock.Controller
	recorder *MockCursorStoreMockRecorder
}

// MockCursorStoreMockRecorder is the mock recorder for MockCursorStore.
type MockCursorStoreMockRecorder struct {
	mock *MockCursorStore
}

// NewMockCursorStore creates a new mock instance.
func NewMockCursorStore(ctrl *gomock.Controller) *MockCursorStore {
	mock := &MockCursorStore{ctrl: ctrl}
	mock.recorder = &MockCursorStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCursorStore) EXPECT() *MockCursorStoreMockRecorder {
	return m.recorder
}

// CreateCursor mocks base method.
func (m *MockCursorStore) CreateCursor(arg0 context.Context, arg1 string, arg2 uint) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCursor", arg0, arg1, arg2)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCursor indicates an expected call of CreateCursor.
func (mr *MockCursorStoreMockRecorder) CreateCursor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCursor", reflect.TypeOf((*MockCursorStore)(nil).CreateCursor), arg0, arg1, arg2)
}

// LookupCursorByID mocks base method.
func (m *MockCursorStore) LookupCursorByID(arg0 context.Context, arg1 string) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCursorByID", arg0, arg1)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCursorByID indicates an expected call of LookupCursorByID.
func (mr *MockCursorStoreMockRecorder) LookupCursorByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCursorByID", reflect.TypeOf((*MockCursorStore)(nil).LookupCursorByID), arg0, arg1)
}

// LookupCursorByName mocks base method.
func (m *MockCursorStore) LookupCursorByName(arg0 context.Context, arg1 string) (flux.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCursorByName", arg0, arg1)
	ret0, _ := ret[0].(flux.Cursor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCursorByName indicates an expected call of LookupCursorByName.
func (mr *MockCursorStoreMockRecorder) LookupCursorByName(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCursorByName", reflect.TypeOf((*MockCursorStore)(nil).LookupCursorByName), arg0, arg1)
}

// UpdateCursor mocks base method.
func (m *MockCursorStore) UpdateCursor(arg0 context.Context, arg1 string, arg2 uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCursor", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCursor indicates an expected call of UpdateCursor.
func (mr *MockCursorStoreMockRecorder) UpdateCursor(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCursor", reflect.TypeOf((*MockCursorStore)(nil).UpdateCursor), arg0, arg1, arg2)
}
