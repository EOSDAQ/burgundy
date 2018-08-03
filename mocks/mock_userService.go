// Code generated by MockGen. DO NOT EDIT.
// Source: service/userService.go

// Package mock_service is a generated GoMock package.
package mocks

import (
	models "burgundy/models"
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserService is a mock of UserService interface
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// GetByID mocks base method
func (m *MockUserService) GetByID(ctx context.Context, accountName string) (*models.User, error) {
	ret := m.ctrl.Call(m, "GetByID", ctx, accountName)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUserServiceMockRecorder) GetByID(ctx, accountName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUserService)(nil).GetByID), ctx, accountName)
}

// Store mocks base method
func (m *MockUserService) Store(ctx context.Context, user *models.User) (*models.User, error) {
	ret := m.ctrl.Call(m, "Store", ctx, user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store
func (mr *MockUserServiceMockRecorder) Store(ctx, user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockUserService)(nil).Store), ctx, user)
}

// Delete mocks base method
func (m *MockUserService) Delete(ctx context.Context, accountName string) (bool, error) {
	ret := m.ctrl.Call(m, "Delete", ctx, accountName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete
func (mr *MockUserServiceMockRecorder) Delete(ctx, accountName interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUserService)(nil).Delete), ctx, accountName)
}
