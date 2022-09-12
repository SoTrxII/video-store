// Source: github.com/dapr/go-sdk/client (interfaces: Client)
// Package mock_client is a generated GoMock package.
package mock_client

import (
	context "context"
	reflect "reflect"

	runtime "github.com/dapr/dapr/pkg/proto/runtime/v1"
	actor "github.com/dapr/go-sdk/actor"
	config "github.com/dapr/go-sdk/actor/config"
	client "github.com/dapr/go-sdk/client"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockClient)(nil).Close))
}

// DeleteBulkState mocks base method.
func (m *MockClient) DeleteBulkState(arg0 context.Context, arg1 string, arg2 []string, arg3 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBulkState", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBulkState indicates an expected call of DeleteBulkState.
func (mr *MockClientMockRecorder) DeleteBulkState(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBulkState", reflect.TypeOf((*MockClient)(nil).DeleteBulkState), arg0, arg1, arg2, arg3)
}

// DeleteBulkStateItems mocks base method.
func (m *MockClient) DeleteBulkStateItems(arg0 context.Context, arg1 string, arg2 []*client.DeleteStateItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBulkStateItems", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBulkStateItems indicates an expected call of DeleteBulkStateItems.
func (mr *MockClientMockRecorder) DeleteBulkStateItems(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBulkStateItems", reflect.TypeOf((*MockClient)(nil).DeleteBulkStateItems), arg0, arg1, arg2)
}

// DeleteState mocks base method.
func (m *MockClient) DeleteState(arg0 context.Context, arg1, arg2 string, arg3 map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteState", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteState indicates an expected call of DeleteState.
func (mr *MockClientMockRecorder) DeleteState(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteState", reflect.TypeOf((*MockClient)(nil).DeleteState), arg0, arg1, arg2, arg3)
}

// DeleteStateWithETag mocks base method.
func (m *MockClient) DeleteStateWithETag(arg0 context.Context, arg1, arg2 string, arg3 *client.ETag, arg4 map[string]string, arg5 *client.StateOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStateWithETag", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStateWithETag indicates an expected call of DeleteStateWithETag.
func (mr *MockClientMockRecorder) DeleteStateWithETag(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStateWithETag", reflect.TypeOf((*MockClient)(nil).DeleteStateWithETag), arg0, arg1, arg2, arg3, arg4, arg5)
}

// ExecuteStateTransaction mocks base method.
func (m *MockClient) ExecuteStateTransaction(arg0 context.Context, arg1 string, arg2 map[string]string, arg3 []*client.StateOperation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExecuteStateTransaction", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExecuteStateTransaction indicates an expected call of ExecuteStateTransaction.
func (mr *MockClientMockRecorder) ExecuteStateTransaction(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteStateTransaction", reflect.TypeOf((*MockClient)(nil).ExecuteStateTransaction), arg0, arg1, arg2, arg3)
}

// GetActorState mocks base method.
func (m *MockClient) GetActorState(arg0 context.Context, arg1 *client.GetActorStateRequest) (*client.GetActorStateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActorState", arg0, arg1)
	ret0, _ := ret[0].(*client.GetActorStateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActorState indicates an expected call of GetActorState.
func (mr *MockClientMockRecorder) GetActorState(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActorState", reflect.TypeOf((*MockClient)(nil).GetActorState), arg0, arg1)
}

// GetBulkSecret mocks base method.
func (m *MockClient) GetBulkSecret(arg0 context.Context, arg1 string, arg2 map[string]string) (map[string]map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBulkSecret", arg0, arg1, arg2)
	ret0, _ := ret[0].(map[string]map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBulkSecret indicates an expected call of GetBulkSecret.
func (mr *MockClientMockRecorder) GetBulkSecret(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBulkSecret", reflect.TypeOf((*MockClient)(nil).GetBulkSecret), arg0, arg1, arg2)
}

// GetBulkState mocks base method.
func (m *MockClient) GetBulkState(arg0 context.Context, arg1 string, arg2 []string, arg3 map[string]string, arg4 int32) ([]*client.BulkStateItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBulkState", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*client.BulkStateItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBulkState indicates an expected call of GetBulkState.
func (mr *MockClientMockRecorder) GetBulkState(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBulkState", reflect.TypeOf((*MockClient)(nil).GetBulkState), arg0, arg1, arg2, arg3, arg4)
}

// GetConfigurationItem mocks base method.
func (m *MockClient) GetConfigurationItem(arg0 context.Context, arg1, arg2 string, arg3 ...client.ConfigurationOpt) (*client.ConfigurationItem, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfigurationItem", varargs...)
	ret0, _ := ret[0].(*client.ConfigurationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigurationItem indicates an expected call of GetConfigurationItem.
func (mr *MockClientMockRecorder) GetConfigurationItem(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigurationItem", reflect.TypeOf((*MockClient)(nil).GetConfigurationItem), varargs...)
}

// GetConfigurationItems mocks base method.
func (m *MockClient) GetConfigurationItems(arg0 context.Context, arg1 string, arg2 []string, arg3 ...client.ConfigurationOpt) ([]*client.ConfigurationItem, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetConfigurationItems", varargs...)
	ret0, _ := ret[0].([]*client.ConfigurationItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfigurationItems indicates an expected call of GetConfigurationItems.
func (mr *MockClientMockRecorder) GetConfigurationItems(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfigurationItems", reflect.TypeOf((*MockClient)(nil).GetConfigurationItems), varargs...)
}

// GetSecret mocks base method.
func (m *MockClient) GetSecret(arg0 context.Context, arg1, arg2 string, arg3 map[string]string) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockClientMockRecorder) GetSecret(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockClient)(nil).GetSecret), arg0, arg1, arg2, arg3)
}

// GetState mocks base method.
func (m *MockClient) GetState(arg0 context.Context, arg1, arg2 string, arg3 map[string]string) (*client.StateItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*client.StateItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockClientMockRecorder) GetState(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockClient)(nil).GetState), arg0, arg1, arg2, arg3)
}

// GetStateWithConsistency mocks base method.
func (m *MockClient) GetStateWithConsistency(arg0 context.Context, arg1, arg2 string, arg3 map[string]string, arg4 client.StateConsistency) (*client.StateItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStateWithConsistency", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(*client.StateItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStateWithConsistency indicates an expected call of GetStateWithConsistency.
func (mr *MockClientMockRecorder) GetStateWithConsistency(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStateWithConsistency", reflect.TypeOf((*MockClient)(nil).GetStateWithConsistency), arg0, arg1, arg2, arg3, arg4)
}

// GrpcClient mocks base method.
func (m *MockClient) GrpcClient() runtime.DaprClient {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GrpcClient")
	ret0, _ := ret[0].(runtime.DaprClient)
	return ret0
}

// GrpcClient indicates an expected call of GrpcClient.
func (mr *MockClientMockRecorder) GrpcClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GrpcClient", reflect.TypeOf((*MockClient)(nil).GrpcClient))
}

// ImplActorClientStub mocks base method.
func (m *MockClient) ImplActorClientStub(arg0 actor.Client, arg1 ...config.Option) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "ImplActorClientStub", varargs...)
}

// ImplActorClientStub indicates an expected call of ImplActorClientStub.
func (mr *MockClientMockRecorder) ImplActorClientStub(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ImplActorClientStub", reflect.TypeOf((*MockClient)(nil).ImplActorClientStub), varargs...)
}

// InvokeActor mocks base method.
func (m *MockClient) InvokeActor(arg0 context.Context, arg1 *client.InvokeActorRequest) (*client.InvokeActorResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeActor", arg0, arg1)
	ret0, _ := ret[0].(*client.InvokeActorResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeActor indicates an expected call of InvokeActor.
func (mr *MockClientMockRecorder) InvokeActor(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeActor", reflect.TypeOf((*MockClient)(nil).InvokeActor), arg0, arg1)
}

// InvokeBinding mocks base method.
func (m *MockClient) InvokeBinding(arg0 context.Context, arg1 *client.InvokeBindingRequest) (*client.BindingEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeBinding", arg0, arg1)
	ret0, _ := ret[0].(*client.BindingEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeBinding indicates an expected call of InvokeBinding.
func (mr *MockClientMockRecorder) InvokeBinding(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeBinding", reflect.TypeOf((*MockClient)(nil).InvokeBinding), arg0, arg1)
}

// InvokeMethod mocks base method.
func (m *MockClient) InvokeMethod(arg0 context.Context, arg1, arg2, arg3 string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeMethod", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeMethod indicates an expected call of InvokeMethod.
func (mr *MockClientMockRecorder) InvokeMethod(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeMethod", reflect.TypeOf((*MockClient)(nil).InvokeMethod), arg0, arg1, arg2, arg3)
}

// InvokeMethodWithContent mocks base method.
func (m *MockClient) InvokeMethodWithContent(arg0 context.Context, arg1, arg2, arg3 string, arg4 *client.DataContent) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeMethodWithContent", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeMethodWithContent indicates an expected call of InvokeMethodWithContent.
func (mr *MockClientMockRecorder) InvokeMethodWithContent(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeMethodWithContent", reflect.TypeOf((*MockClient)(nil).InvokeMethodWithContent), arg0, arg1, arg2, arg3, arg4)
}

// InvokeMethodWithCustomContent mocks base method.
func (m *MockClient) InvokeMethodWithCustomContent(arg0 context.Context, arg1, arg2, arg3, arg4 string, arg5 interface{}) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeMethodWithCustomContent", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InvokeMethodWithCustomContent indicates an expected call of InvokeMethodWithCustomContent.
func (mr *MockClientMockRecorder) InvokeMethodWithCustomContent(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeMethodWithCustomContent", reflect.TypeOf((*MockClient)(nil).InvokeMethodWithCustomContent), arg0, arg1, arg2, arg3, arg4, arg5)
}

// InvokeOutputBinding mocks base method.
func (m *MockClient) InvokeOutputBinding(arg0 context.Context, arg1 *client.InvokeBindingRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InvokeOutputBinding", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InvokeOutputBinding indicates an expected call of InvokeOutputBinding.
func (mr *MockClientMockRecorder) InvokeOutputBinding(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InvokeOutputBinding", reflect.TypeOf((*MockClient)(nil).InvokeOutputBinding), arg0, arg1)
}

// PublishEvent mocks base method.
func (m *MockClient) PublishEvent(arg0 context.Context, arg1, arg2 string, arg3 interface{}, arg4 ...client.PublishEventOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PublishEvent", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEvent indicates an expected call of PublishEvent.
func (mr *MockClientMockRecorder) PublishEvent(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEvent", reflect.TypeOf((*MockClient)(nil).PublishEvent), varargs...)
}

// PublishEventfromCustomContent mocks base method.
func (m *MockClient) PublishEventfromCustomContent(arg0 context.Context, arg1, arg2 string, arg3 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishEventfromCustomContent", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishEventfromCustomContent indicates an expected call of PublishEventfromCustomContent.
func (mr *MockClientMockRecorder) PublishEventfromCustomContent(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishEventfromCustomContent", reflect.TypeOf((*MockClient)(nil).PublishEventfromCustomContent), arg0, arg1, arg2, arg3)
}

// QueryStateAlpha1 mocks base method.
func (m *MockClient) QueryStateAlpha1(arg0 context.Context, arg1, arg2 string, arg3 map[string]string) (*client.QueryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryStateAlpha1", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*client.QueryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryStateAlpha1 indicates an expected call of QueryStateAlpha1.
func (mr *MockClientMockRecorder) QueryStateAlpha1(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryStateAlpha1", reflect.TypeOf((*MockClient)(nil).QueryStateAlpha1), arg0, arg1, arg2, arg3)
}

// RegisterActorReminder mocks base method.
func (m *MockClient) RegisterActorReminder(arg0 context.Context, arg1 *client.RegisterActorReminderRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterActorReminder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterActorReminder indicates an expected call of RegisterActorReminder.
func (mr *MockClientMockRecorder) RegisterActorReminder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterActorReminder", reflect.TypeOf((*MockClient)(nil).RegisterActorReminder), arg0, arg1)
}

// RegisterActorTimer mocks base method.
func (m *MockClient) RegisterActorTimer(arg0 context.Context, arg1 *client.RegisterActorTimerRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterActorTimer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterActorTimer indicates an expected call of RegisterActorTimer.
func (mr *MockClientMockRecorder) RegisterActorTimer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterActorTimer", reflect.TypeOf((*MockClient)(nil).RegisterActorTimer), arg0, arg1)
}

// RenameActorReminder mocks base method.
func (m *MockClient) RenameActorReminder(arg0 context.Context, arg1 *client.RenameActorReminderRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RenameActorReminder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameActorReminder indicates an expected call of RenameActorReminder.
func (mr *MockClientMockRecorder) RenameActorReminder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RenameActorReminder", reflect.TypeOf((*MockClient)(nil).RenameActorReminder), arg0, arg1)
}

// SaveBulkState mocks base method.
func (m *MockClient) SaveBulkState(arg0 context.Context, arg1 string, arg2 ...*client.SetStateItem) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveBulkState", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBulkState indicates an expected call of SaveBulkState.
func (mr *MockClientMockRecorder) SaveBulkState(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBulkState", reflect.TypeOf((*MockClient)(nil).SaveBulkState), varargs...)
}

// SaveState mocks base method.
func (m *MockClient) SaveState(arg0 context.Context, arg1, arg2 string, arg3 []byte, arg4 map[string]string, arg5 ...client.StateOption) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3, arg4}
	for _, a := range arg5 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveState", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveState indicates an expected call of SaveState.
func (mr *MockClientMockRecorder) SaveState(arg0, arg1, arg2, arg3, arg4 interface{}, arg5 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3, arg4}, arg5...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveState", reflect.TypeOf((*MockClient)(nil).SaveState), varargs...)
}

// SaveStateTransactionally mocks base method.
func (m *MockClient) SaveStateTransactionally(arg0 context.Context, arg1, arg2 string, arg3 []*client.ActorStateOperation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveStateTransactionally", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveStateTransactionally indicates an expected call of SaveStateTransactionally.
func (mr *MockClientMockRecorder) SaveStateTransactionally(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveStateTransactionally", reflect.TypeOf((*MockClient)(nil).SaveStateTransactionally), arg0, arg1, arg2, arg3)
}

// Shutdown mocks base method.
func (m *MockClient) Shutdown(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Shutdown", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockClientMockRecorder) Shutdown(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockClient)(nil).Shutdown), arg0)
}

// SubscribeConfigurationItems mocks base method.
func (m *MockClient) SubscribeConfigurationItems(arg0 context.Context, arg1 string, arg2 []string, arg3 client.ConfigurationHandleFunction, arg4 ...client.ConfigurationOpt) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2, arg3}
	for _, a := range arg4 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SubscribeConfigurationItems", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SubscribeConfigurationItems indicates an expected call of SubscribeConfigurationItems.
func (mr *MockClientMockRecorder) SubscribeConfigurationItems(arg0, arg1, arg2, arg3 interface{}, arg4 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2, arg3}, arg4...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeConfigurationItems", reflect.TypeOf((*MockClient)(nil).SubscribeConfigurationItems), varargs...)
}

// TryLockAlpha1 mocks base method.
func (m *MockClient) TryLockAlpha1(arg0 context.Context, arg1 string, arg2 *client.LockRequest) (*client.LockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TryLockAlpha1", arg0, arg1, arg2)
	ret0, _ := ret[0].(*client.LockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TryLockAlpha1 indicates an expected call of TryLockAlpha1.
func (mr *MockClientMockRecorder) TryLockAlpha1(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TryLockAlpha1", reflect.TypeOf((*MockClient)(nil).TryLockAlpha1), arg0, arg1, arg2)
}

// UnlockAlpha1 mocks base method.
func (m *MockClient) UnlockAlpha1(arg0 context.Context, arg1 string, arg2 *client.UnlockRequest) (*client.UnlockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlockAlpha1", arg0, arg1, arg2)
	ret0, _ := ret[0].(*client.UnlockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnlockAlpha1 indicates an expected call of UnlockAlpha1.
func (mr *MockClientMockRecorder) UnlockAlpha1(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlockAlpha1", reflect.TypeOf((*MockClient)(nil).UnlockAlpha1), arg0, arg1, arg2)
}

// UnregisterActorReminder mocks base method.
func (m *MockClient) UnregisterActorReminder(arg0 context.Context, arg1 *client.UnregisterActorReminderRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterActorReminder", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterActorReminder indicates an expected call of UnregisterActorReminder.
func (mr *MockClientMockRecorder) UnregisterActorReminder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterActorReminder", reflect.TypeOf((*MockClient)(nil).UnregisterActorReminder), arg0, arg1)
}

// UnregisterActorTimer mocks base method.
func (m *MockClient) UnregisterActorTimer(arg0 context.Context, arg1 *client.UnregisterActorTimerRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnregisterActorTimer", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnregisterActorTimer indicates an expected call of UnregisterActorTimer.
func (mr *MockClientMockRecorder) UnregisterActorTimer(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnregisterActorTimer", reflect.TypeOf((*MockClient)(nil).UnregisterActorTimer), arg0, arg1)
}

// UnsubscribeConfigurationItems mocks base method.
func (m *MockClient) UnsubscribeConfigurationItems(arg0 context.Context, arg1, arg2 string, arg3 ...client.ConfigurationOpt) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1, arg2}
	for _, a := range arg3 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UnsubscribeConfigurationItems", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UnsubscribeConfigurationItems indicates an expected call of UnsubscribeConfigurationItems.
func (mr *MockClientMockRecorder) UnsubscribeConfigurationItems(arg0, arg1, arg2 interface{}, arg3 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1, arg2}, arg3...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnsubscribeConfigurationItems", reflect.TypeOf((*MockClient)(nil).UnsubscribeConfigurationItems), varargs...)
}

// WithAuthToken mocks base method.
func (m *MockClient) WithAuthToken(arg0 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "WithAuthToken", arg0)
}

// WithAuthToken indicates an expected call of WithAuthToken.
func (mr *MockClientMockRecorder) WithAuthToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithAuthToken", reflect.TypeOf((*MockClient)(nil).WithAuthToken), arg0)
}

// WithTraceID mocks base method.
func (m *MockClient) WithTraceID(arg0 context.Context, arg1 string) context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTraceID", arg0, arg1)
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// WithTraceID indicates an expected call of WithTraceID.
func (mr *MockClientMockRecorder) WithTraceID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTraceID", reflect.TypeOf((*MockClient)(nil).WithTraceID), arg0, arg1)
}
