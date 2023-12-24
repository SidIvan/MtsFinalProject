// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/driver_service.go

// Package mock_service is a generated GoMock package.
package mock

import (
	context "context"
	model "driver-service/internal/model"
	service "driver-service/internal/svc"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDriverRepo is a mock of DriverRepo interface.
type MockDriverRepo struct {
	ctrl     *gomock.Controller
	recorder *MockDriverRepoMockRecorder
}

// MockDriverRepoMockRecorder is the mock recorder for MockDriverRepo.
type MockDriverRepoMockRecorder struct {
	mock *MockDriverRepo
}

// NewMockDriverRepo creates a new mock instance.
func NewMockDriverRepo(ctrl *gomock.Controller) *MockDriverRepo {
	mock := &MockDriverRepo{ctrl: ctrl}
	mock.recorder = &MockDriverRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDriverRepo) EXPECT() *MockDriverRepoMockRecorder {
	return m.recorder
}

// CreateTrip mocks base method.
func (m *MockDriverRepo) CreateTrip(arg0 context.Context, arg1 *model.Trip) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrip", arg0, arg1)
	ret0, _ := ret[0].(string)
	return ret0
}

// CreateTrip indicates an expected call of CreateTrip.
func (mr *MockDriverRepoMockRecorder) CreateTrip(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrip", reflect.TypeOf((*MockDriverRepo)(nil).CreateTrip), arg0, arg1)
}

// GetTrip mocks base method.
func (m *MockDriverRepo) GetTrip(arg0 context.Context, arg1, arg2 string) *model.Trip {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Trip)
	return ret0
}

// GetTrip indicates an expected call of GetTrip.
func (mr *MockDriverRepoMockRecorder) GetTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrip", reflect.TypeOf((*MockDriverRepo)(nil).GetTrip), arg0, arg1, arg2)
}

// GetTrips mocks base method.
func (m *MockDriverRepo) GetTrips(arg0 context.Context, arg1 string) []model.Trip {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrips", arg0, arg1)
	ret0, _ := ret[0].([]model.Trip)
	return ret0
}

// GetTrips indicates an expected call of GetTrips.
func (mr *MockDriverRepoMockRecorder) GetTrips(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrips", reflect.TypeOf((*MockDriverRepo)(nil).GetTrips), arg0, arg1)
}

// PutTripStatus mocks base method.
func (m *MockDriverRepo) PutTripStatus(arg0 context.Context, arg1, arg2 string, arg3 model.TripStatus) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutTripStatus", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(bool)
	return ret0
}

// PutTripStatus indicates an expected call of PutTripStatus.
func (mr *MockDriverRepoMockRecorder) PutTripStatus(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutTripStatus", reflect.TypeOf((*MockDriverRepo)(nil).PutTripStatus), arg0, arg1, arg2, arg3)
}

// SaveCancelReason mocks base method.
func (m *MockDriverRepo) SaveCancelReason(arg0 context.Context, arg1, arg2, arg3 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SaveCancelReason", arg0, arg1, arg2, arg3)
}

// SaveCancelReason indicates an expected call of SaveCancelReason.
func (mr *MockDriverRepoMockRecorder) SaveCancelReason(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCancelReason", reflect.TypeOf((*MockDriverRepo)(nil).SaveCancelReason), arg0, arg1, arg2, arg3)
}

// MockLocationClient is a mock of LocationClient interface.
type MockLocationClient struct {
	ctrl     *gomock.Controller
	recorder *MockLocationClientMockRecorder
}

// MockLocationClientMockRecorder is the mock recorder for MockLocationClient.
type MockLocationClientMockRecorder struct {
	mock *MockLocationClient
}

// NewMockLocationClient creates a new mock instance.
func NewMockLocationClient(ctrl *gomock.Controller) *MockLocationClient {
	mock := &MockLocationClient{ctrl: ctrl}
	mock.recorder = &MockLocationClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLocationClient) EXPECT() *MockLocationClientMockRecorder {
	return m.recorder
}

// GetDrivers mocks base method.
func (m *MockLocationClient) GetDrivers(arg0 context.Context, arg1 *service.GetDriversPayload) []model.Driver {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDrivers", arg0, arg1)
	ret0, _ := ret[0].([]model.Driver)
	return ret0
}

// GetDrivers indicates an expected call of GetDrivers.
func (mr *MockLocationClientMockRecorder) GetDrivers(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDrivers", reflect.TypeOf((*MockLocationClient)(nil).GetDrivers), arg0, arg1)
}

// MockEventProducer is a mock of EventProducer interface.
type MockEventProducer struct {
	ctrl     *gomock.Controller
	recorder *MockEventProducerMockRecorder
}

// MockEventProducerMockRecorder is the mock recorder for MockEventProducer.
type MockEventProducerMockRecorder struct {
	mock *MockEventProducer
}

// NewMockEventProducer creates a new mock instance.
func NewMockEventProducer(ctrl *gomock.Controller) *MockEventProducer {
	mock := &MockEventProducer{ctrl: ctrl}
	mock.recorder = &MockEventProducerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventProducer) EXPECT() *MockEventProducerMockRecorder {
	return m.recorder
}

// Free mocks base method.
func (m *MockEventProducer) Free() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Free")
}

// Free indicates an expected call of Free.
func (mr *MockEventProducerMockRecorder) Free() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Free", reflect.TypeOf((*MockEventProducer)(nil).Free))
}

// SendTripEvent mocks base method.
func (m *MockEventProducer) SendTripEvent(arg0 context.Context, arg1 *service.TripMessagePayload) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SendTripEvent", arg0, arg1)
}

// SendTripEvent indicates an expected call of SendTripEvent.
func (mr *MockEventProducerMockRecorder) SendTripEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendTripEvent", reflect.TypeOf((*MockEventProducer)(nil).SendTripEvent), arg0, arg1)
}

// MockDriverService is a mock of DriverService interface.
type MockDriverService struct {
	ctrl     *gomock.Controller
	recorder *MockDriverServiceMockRecorder
}

// MockDriverServiceMockRecorder is the mock recorder for MockDriverService.
type MockDriverServiceMockRecorder struct {
	mock *MockDriverService
}

// NewMockDriverService creates a new mock instance.
func NewMockDriverService(ctrl *gomock.Controller) *MockDriverService {
	mock := &MockDriverService{ctrl: ctrl}
	mock.recorder = &MockDriverServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDriverService) EXPECT() *MockDriverServiceMockRecorder {
	return m.recorder
}

// AcceptTrip mocks base method.
func (m *MockDriverService) AcceptTrip(arg0 context.Context, arg1, arg2 string) service.ChangeTripStatusStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AcceptTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(service.ChangeTripStatusStatus)
	return ret0
}

// AcceptTrip indicates an expected call of AcceptTrip.
func (mr *MockDriverServiceMockRecorder) AcceptTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AcceptTrip", reflect.TypeOf((*MockDriverService)(nil).AcceptTrip), arg0, arg1, arg2)
}

// CancelTrip mocks base method.
func (m *MockDriverService) CancelTrip(arg0 context.Context, arg1, arg2, arg3 string) service.ChangeTripStatusStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CancelTrip", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(service.ChangeTripStatusStatus)
	return ret0
}

// CancelTrip indicates an expected call of CancelTrip.
func (mr *MockDriverServiceMockRecorder) CancelTrip(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CancelTrip", reflect.TypeOf((*MockDriverService)(nil).CancelTrip), arg0, arg1, arg2, arg3)
}

// CloseEventWriter mocks base method.
func (m *MockDriverService) CloseEventWriter() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseEventWriter")
}

// CloseEventWriter indicates an expected call of CloseEventWriter.
func (mr *MockDriverServiceMockRecorder) CloseEventWriter() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseEventWriter", reflect.TypeOf((*MockDriverService)(nil).CloseEventWriter))
}

// CreateTrip mocks base method.
func (m *MockDriverService) CreateTrip(arg0 context.Context, arg1 model.Trip) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateTrip", arg0, arg1)
}

// CreateTrip indicates an expected call of CreateTrip.
func (mr *MockDriverServiceMockRecorder) CreateTrip(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrip", reflect.TypeOf((*MockDriverService)(nil).CreateTrip), arg0, arg1)
}

// EndTrip mocks base method.
func (m *MockDriverService) EndTrip(arg0 context.Context, arg1, arg2 string) service.ChangeTripStatusStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EndTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(service.ChangeTripStatusStatus)
	return ret0
}

// EndTrip indicates an expected call of EndTrip.
func (mr *MockDriverServiceMockRecorder) EndTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EndTrip", reflect.TypeOf((*MockDriverService)(nil).EndTrip), arg0, arg1, arg2)
}

// GetTrip mocks base method.
func (m *MockDriverService) GetTrip(arg0 context.Context, arg1, arg2 string) *model.Trip {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.Trip)
	return ret0
}

// GetTrip indicates an expected call of GetTrip.
func (mr *MockDriverServiceMockRecorder) GetTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrip", reflect.TypeOf((*MockDriverService)(nil).GetTrip), arg0, arg1, arg2)
}

// GetTrips mocks base method.
func (m *MockDriverService) GetTrips(arg0 context.Context, arg1 string) []model.Trip {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTrips", arg0, arg1)
	ret0, _ := ret[0].([]model.Trip)
	return ret0
}

// GetTrips indicates an expected call of GetTrips.
func (mr *MockDriverServiceMockRecorder) GetTrips(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTrips", reflect.TypeOf((*MockDriverService)(nil).GetTrips), arg0, arg1)
}

// StartTrip mocks base method.
func (m *MockDriverService) StartTrip(arg0 context.Context, arg1, arg2 string) service.ChangeTripStatusStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTrip", arg0, arg1, arg2)
	ret0, _ := ret[0].(service.ChangeTripStatusStatus)
	return ret0
}

// StartTrip indicates an expected call of StartTrip.
func (mr *MockDriverServiceMockRecorder) StartTrip(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTrip", reflect.TypeOf((*MockDriverService)(nil).StartTrip), arg0, arg1, arg2)
}