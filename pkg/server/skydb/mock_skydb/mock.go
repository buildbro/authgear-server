// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/skygeario/skygear-server/pkg/server/skydb (interfaces: Conn,Database)

package mock_skydb

import (
	gomock "github.com/golang/mock/gomock"
	skydb "github.com/skygeario/skygear-server/pkg/server/skydb"
	time "time"
)

// MockConn is a mock of Conn interface
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockConn) EXPECT() *MockConnMockRecorder {
	return _m.recorder
}

// AddRelation mocks base method
func (_m *MockConn) AddRelation(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "AddRelation", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddRelation indicates an expected call of AddRelation
func (_mr *MockConnMockRecorder) AddRelation(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "AddRelation", arg0, arg1, arg2)
}

// Close mocks base method
func (_m *MockConn) Close() error {
	ret := _m.ctrl.Call(_m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (_mr *MockConnMockRecorder) Close() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Close")
}

// CreateUser mocks base method
func (_m *MockConn) CreateUser(_param0 *skydb.UserInfo) error {
	ret := _m.ctrl.Call(_m, "CreateUser", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateUser indicates an expected call of CreateUser
func (_mr *MockConnMockRecorder) CreateUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "CreateUser", arg0)
}

// DeleteDevice mocks base method
func (_m *MockConn) DeleteDevice(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteDevice", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevice indicates an expected call of DeleteDevice
func (_mr *MockConnMockRecorder) DeleteDevice(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteDevice", arg0)
}

// DeleteDevicesByToken mocks base method
func (_m *MockConn) DeleteDevicesByToken(_param0 string, _param1 time.Time) error {
	ret := _m.ctrl.Call(_m, "DeleteDevicesByToken", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteDevicesByToken indicates an expected call of DeleteDevicesByToken
func (_mr *MockConnMockRecorder) DeleteDevicesByToken(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteDevicesByToken", arg0, arg1)
}

// DeleteEmptyDevicesByTime mocks base method
func (_m *MockConn) DeleteEmptyDevicesByTime(_param0 time.Time) error {
	ret := _m.ctrl.Call(_m, "DeleteEmptyDevicesByTime", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEmptyDevicesByTime indicates an expected call of DeleteEmptyDevicesByTime
func (_mr *MockConnMockRecorder) DeleteEmptyDevicesByTime(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteEmptyDevicesByTime", arg0)
}

// DeleteUser mocks base method
func (_m *MockConn) DeleteUser(_param0 string) error {
	ret := _m.ctrl.Call(_m, "DeleteUser", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (_mr *MockConnMockRecorder) DeleteUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteUser", arg0)
}

// GetAdminRoles mocks base method
func (_m *MockConn) GetAdminRoles() ([]string, error) {
	ret := _m.ctrl.Call(_m, "GetAdminRoles")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminRoles indicates an expected call of GetAdminRoles
func (_mr *MockConnMockRecorder) GetAdminRoles() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAdminRoles")
}

// GetAsset mocks base method
func (_m *MockConn) GetAsset(_param0 string, _param1 *skydb.Asset) error {
	ret := _m.ctrl.Call(_m, "GetAsset", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAsset indicates an expected call of GetAsset
func (_mr *MockConnMockRecorder) GetAsset(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAsset", arg0, arg1)
}

// GetAssets mocks base method
func (_m *MockConn) GetAssets(_param0 []string) ([]skydb.Asset, error) {
	ret := _m.ctrl.Call(_m, "GetAssets", _param0)
	ret0, _ := ret[0].([]skydb.Asset)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssets indicates an expected call of GetAssets
func (_mr *MockConnMockRecorder) GetAssets(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAssets", arg0)
}

// GetDefaultRoles mocks base method
func (_m *MockConn) GetDefaultRoles() ([]string, error) {
	ret := _m.ctrl.Call(_m, "GetDefaultRoles")
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDefaultRoles indicates an expected call of GetDefaultRoles
func (_mr *MockConnMockRecorder) GetDefaultRoles() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDefaultRoles")
}

// GetDevice mocks base method
func (_m *MockConn) GetDevice(_param0 string, _param1 *skydb.Device) error {
	ret := _m.ctrl.Call(_m, "GetDevice", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetDevice indicates an expected call of GetDevice
func (_mr *MockConnMockRecorder) GetDevice(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetDevice", arg0, arg1)
}

// GetRecordAccess mocks base method
func (_m *MockConn) GetRecordAccess(_param0 string) (skydb.RecordACL, error) {
	ret := _m.ctrl.Call(_m, "GetRecordAccess", _param0)
	ret0, _ := ret[0].(skydb.RecordACL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordAccess indicates an expected call of GetRecordAccess
func (_mr *MockConnMockRecorder) GetRecordAccess(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRecordAccess", arg0)
}

// GetRecordDefaultAccess mocks base method
func (_m *MockConn) GetRecordDefaultAccess(_param0 string) (skydb.RecordACL, error) {
	ret := _m.ctrl.Call(_m, "GetRecordDefaultAccess", _param0)
	ret0, _ := ret[0].(skydb.RecordACL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordDefaultAccess indicates an expected call of GetRecordDefaultAccess
func (_mr *MockConnMockRecorder) GetRecordDefaultAccess(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRecordDefaultAccess", arg0)
}

// GetUser mocks base method
func (_m *MockConn) GetUser(_param0 string, _param1 *skydb.UserInfo) error {
	ret := _m.ctrl.Call(_m, "GetUser", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUser indicates an expected call of GetUser
func (_mr *MockConnMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUser", arg0, arg1)
}

// GetUserByPrincipalID mocks base method
func (_m *MockConn) GetUserByPrincipalID(_param0 string, _param1 *skydb.UserInfo) error {
	ret := _m.ctrl.Call(_m, "GetUserByPrincipalID", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserByPrincipalID indicates an expected call of GetUserByPrincipalID
func (_mr *MockConnMockRecorder) GetUserByPrincipalID(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUserByPrincipalID", arg0, arg1)
}

// GetUserByUsernameEmail mocks base method
func (_m *MockConn) GetUserByUsernameEmail(_param0 string, _param1 string, _param2 *skydb.UserInfo) error {
	ret := _m.ctrl.Call(_m, "GetUserByUsernameEmail", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetUserByUsernameEmail indicates an expected call of GetUserByUsernameEmail
func (_mr *MockConnMockRecorder) GetUserByUsernameEmail(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetUserByUsernameEmail", arg0, arg1, arg2)
}

// PrivateDB mocks base method
func (_m *MockConn) PrivateDB(_param0 string) skydb.Database {
	ret := _m.ctrl.Call(_m, "PrivateDB", _param0)
	ret0, _ := ret[0].(skydb.Database)
	return ret0
}

// PrivateDB indicates an expected call of PrivateDB
func (_mr *MockConnMockRecorder) PrivateDB(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PrivateDB", arg0)
}

// PublicDB mocks base method
func (_m *MockConn) PublicDB() skydb.Database {
	ret := _m.ctrl.Call(_m, "PublicDB")
	ret0, _ := ret[0].(skydb.Database)
	return ret0
}

// PublicDB indicates an expected call of PublicDB
func (_mr *MockConnMockRecorder) PublicDB() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "PublicDB")
}

// QueryDevicesByUser mocks base method
func (_m *MockConn) QueryDevicesByUser(_param0 string) ([]skydb.Device, error) {
	ret := _m.ctrl.Call(_m, "QueryDevicesByUser", _param0)
	ret0, _ := ret[0].([]skydb.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryDevicesByUser indicates an expected call of QueryDevicesByUser
func (_mr *MockConnMockRecorder) QueryDevicesByUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryDevicesByUser", arg0)
}

// QueryDevicesByUserAndTopic mocks base method
func (_m *MockConn) QueryDevicesByUserAndTopic(_param0 string, _param1 string) ([]skydb.Device, error) {
	ret := _m.ctrl.Call(_m, "QueryDevicesByUserAndTopic", _param0, _param1)
	ret0, _ := ret[0].([]skydb.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryDevicesByUserAndTopic indicates an expected call of QueryDevicesByUserAndTopic
func (_mr *MockConnMockRecorder) QueryDevicesByUserAndTopic(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryDevicesByUserAndTopic", arg0, arg1)
}

// QueryRelation mocks base method
func (_m *MockConn) QueryRelation(_param0 string, _param1 string, _param2 string, _param3 skydb.QueryConfig) []skydb.UserInfo {
	ret := _m.ctrl.Call(_m, "QueryRelation", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].([]skydb.UserInfo)
	return ret0
}

// QueryRelation indicates an expected call of QueryRelation
func (_mr *MockConnMockRecorder) QueryRelation(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryRelation", arg0, arg1, arg2, arg3)
}

// QueryRelationCount mocks base method
func (_m *MockConn) QueryRelationCount(_param0 string, _param1 string, _param2 string) (uint64, error) {
	ret := _m.ctrl.Call(_m, "QueryRelationCount", _param0, _param1, _param2)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryRelationCount indicates an expected call of QueryRelationCount
func (_mr *MockConnMockRecorder) QueryRelationCount(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryRelationCount", arg0, arg1, arg2)
}

// QueryUser mocks base method
func (_m *MockConn) QueryUser(_param0 []string, _param1 []string) ([]skydb.UserInfo, error) {
	ret := _m.ctrl.Call(_m, "QueryUser", _param0, _param1)
	ret0, _ := ret[0].([]skydb.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryUser indicates an expected call of QueryUser
func (_mr *MockConnMockRecorder) QueryUser(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryUser", arg0, arg1)
}

// RemoveRelation mocks base method
func (_m *MockConn) RemoveRelation(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "RemoveRelation", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRelation indicates an expected call of RemoveRelation
func (_mr *MockConnMockRecorder) RemoveRelation(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoveRelation", arg0, arg1, arg2)
}

// SaveAsset mocks base method
func (_m *MockConn) SaveAsset(_param0 *skydb.Asset) error {
	ret := _m.ctrl.Call(_m, "SaveAsset", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAsset indicates an expected call of SaveAsset
func (_mr *MockConnMockRecorder) SaveAsset(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SaveAsset", arg0)
}

// SaveDevice mocks base method
func (_m *MockConn) SaveDevice(_param0 *skydb.Device) error {
	ret := _m.ctrl.Call(_m, "SaveDevice", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDevice indicates an expected call of SaveDevice
func (_mr *MockConnMockRecorder) SaveDevice(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SaveDevice", arg0)
}

// SetAdminRoles mocks base method
func (_m *MockConn) SetAdminRoles(_param0 []string) error {
	ret := _m.ctrl.Call(_m, "SetAdminRoles", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetAdminRoles indicates an expected call of SetAdminRoles
func (_mr *MockConnMockRecorder) SetAdminRoles(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetAdminRoles", arg0)
}

// SetDefaultRoles mocks base method
func (_m *MockConn) SetDefaultRoles(_param0 []string) error {
	ret := _m.ctrl.Call(_m, "SetDefaultRoles", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetDefaultRoles indicates an expected call of SetDefaultRoles
func (_mr *MockConnMockRecorder) SetDefaultRoles(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetDefaultRoles", arg0)
}

// SetRecordAccess mocks base method
func (_m *MockConn) SetRecordAccess(_param0 string, _param1 skydb.RecordACL) error {
	ret := _m.ctrl.Call(_m, "SetRecordAccess", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecordAccess indicates an expected call of SetRecordAccess
func (_mr *MockConnMockRecorder) SetRecordAccess(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRecordAccess", arg0, arg1)
}

// SetRecordDefaultAccess mocks base method
func (_m *MockConn) SetRecordDefaultAccess(_param0 string, _param1 skydb.RecordACL) error {
	ret := _m.ctrl.Call(_m, "SetRecordDefaultAccess", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetRecordDefaultAccess indicates an expected call of SetRecordDefaultAccess
func (_mr *MockConnMockRecorder) SetRecordDefaultAccess(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SetRecordDefaultAccess", arg0, arg1)
}

// Subscribe mocks base method
func (_m *MockConn) Subscribe(_param0 chan skydb.RecordEvent) error {
	ret := _m.ctrl.Call(_m, "Subscribe", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Subscribe indicates an expected call of Subscribe
func (_mr *MockConnMockRecorder) Subscribe(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Subscribe", arg0)
}

// UnionDB mocks base method
func (_m *MockConn) UnionDB() skydb.Database {
	ret := _m.ctrl.Call(_m, "UnionDB")
	ret0, _ := ret[0].(skydb.Database)
	return ret0
}

// UnionDB indicates an expected call of UnionDB
func (_mr *MockConnMockRecorder) UnionDB() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UnionDB")
}

// UpdateUser mocks base method
func (_m *MockConn) UpdateUser(_param0 *skydb.UserInfo) error {
	ret := _m.ctrl.Call(_m, "UpdateUser", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateUser indicates an expected call of UpdateUser
func (_mr *MockConnMockRecorder) UpdateUser(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateUser", arg0)
}

// MockDatabase is a mock of Database interface
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (_m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return _m.recorder
}

// Conn mocks base method
func (_m *MockDatabase) Conn() skydb.Conn {
	ret := _m.ctrl.Call(_m, "Conn")
	ret0, _ := ret[0].(skydb.Conn)
	return ret0
}

// Conn indicates an expected call of Conn
func (_mr *MockDatabaseMockRecorder) Conn() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Conn")
}

// DatabaseType mocks base method
func (_m *MockDatabase) DatabaseType() skydb.DatabaseType {
	ret := _m.ctrl.Call(_m, "DatabaseType")
	ret0, _ := ret[0].(skydb.DatabaseType)
	return ret0
}

// DatabaseType indicates an expected call of DatabaseType
func (_mr *MockDatabaseMockRecorder) DatabaseType() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DatabaseType")
}

// Delete mocks base method
func (_m *MockDatabase) Delete(_param0 skydb.RecordID) error {
	ret := _m.ctrl.Call(_m, "Delete", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (_mr *MockDatabaseMockRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

// DeleteSchema mocks base method
func (_m *MockDatabase) DeleteSchema(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteSchema", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSchema indicates an expected call of DeleteSchema
func (_mr *MockDatabaseMockRecorder) DeleteSchema(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteSchema", arg0, arg1)
}

// DeleteSubscription mocks base method
func (_m *MockDatabase) DeleteSubscription(_param0 string, _param1 string) error {
	ret := _m.ctrl.Call(_m, "DeleteSubscription", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSubscription indicates an expected call of DeleteSubscription
func (_mr *MockDatabaseMockRecorder) DeleteSubscription(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteSubscription", arg0, arg1)
}

// Extend mocks base method
func (_m *MockDatabase) Extend(_param0 string, _param1 skydb.RecordSchema) (bool, error) {
	ret := _m.ctrl.Call(_m, "Extend", _param0, _param1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Extend indicates an expected call of Extend
func (_mr *MockDatabaseMockRecorder) Extend(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Extend", arg0, arg1)
}

// Get mocks base method
func (_m *MockDatabase) Get(_param0 skydb.RecordID, _param1 *skydb.Record) error {
	ret := _m.ctrl.Call(_m, "Get", _param0, _param1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get
func (_mr *MockDatabaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0, arg1)
}

// GetByIDs mocks base method
func (_m *MockDatabase) GetByIDs(_param0 []skydb.RecordID) (*skydb.Rows, error) {
	ret := _m.ctrl.Call(_m, "GetByIDs", _param0)
	ret0, _ := ret[0].(*skydb.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByIDs indicates an expected call of GetByIDs
func (_mr *MockDatabaseMockRecorder) GetByIDs(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetByIDs", arg0)
}

// GetMatchingSubscriptions mocks base method
func (_m *MockDatabase) GetMatchingSubscriptions(_param0 *skydb.Record) []skydb.Subscription {
	ret := _m.ctrl.Call(_m, "GetMatchingSubscriptions", _param0)
	ret0, _ := ret[0].([]skydb.Subscription)
	return ret0
}

// GetMatchingSubscriptions indicates an expected call of GetMatchingSubscriptions
func (_mr *MockDatabaseMockRecorder) GetMatchingSubscriptions(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetMatchingSubscriptions", arg0)
}

// GetRecordSchemas mocks base method
func (_m *MockDatabase) GetRecordSchemas() (map[string]skydb.RecordSchema, error) {
	ret := _m.ctrl.Call(_m, "GetRecordSchemas")
	ret0, _ := ret[0].(map[string]skydb.RecordSchema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRecordSchemas indicates an expected call of GetRecordSchemas
func (_mr *MockDatabaseMockRecorder) GetRecordSchemas() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetRecordSchemas")
}

// GetSchema mocks base method
func (_m *MockDatabase) GetSchema(_param0 string) (skydb.RecordSchema, error) {
	ret := _m.ctrl.Call(_m, "GetSchema", _param0)
	ret0, _ := ret[0].(skydb.RecordSchema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSchema indicates an expected call of GetSchema
func (_mr *MockDatabaseMockRecorder) GetSchema(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSchema", arg0)
}

// GetSubscription mocks base method
func (_m *MockDatabase) GetSubscription(_param0 string, _param1 string, _param2 *skydb.Subscription) error {
	ret := _m.ctrl.Call(_m, "GetSubscription", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetSubscription indicates an expected call of GetSubscription
func (_mr *MockDatabaseMockRecorder) GetSubscription(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSubscription", arg0, arg1, arg2)
}

// GetSubscriptionsByDeviceID mocks base method
func (_m *MockDatabase) GetSubscriptionsByDeviceID(_param0 string) []skydb.Subscription {
	ret := _m.ctrl.Call(_m, "GetSubscriptionsByDeviceID", _param0)
	ret0, _ := ret[0].([]skydb.Subscription)
	return ret0
}

// GetSubscriptionsByDeviceID indicates an expected call of GetSubscriptionsByDeviceID
func (_mr *MockDatabaseMockRecorder) GetSubscriptionsByDeviceID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetSubscriptionsByDeviceID", arg0)
}

// ID mocks base method
func (_m *MockDatabase) ID() string {
	ret := _m.ctrl.Call(_m, "ID")
	ret0, _ := ret[0].(string)
	return ret0
}

// ID indicates an expected call of ID
func (_mr *MockDatabaseMockRecorder) ID() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "ID")
}

// IsReadOnly mocks base method
func (_m *MockDatabase) IsReadOnly() bool {
	ret := _m.ctrl.Call(_m, "IsReadOnly")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsReadOnly indicates an expected call of IsReadOnly
func (_mr *MockDatabaseMockRecorder) IsReadOnly() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsReadOnly")
}

// Query mocks base method
func (_m *MockDatabase) Query(_param0 *skydb.Query) (*skydb.Rows, error) {
	ret := _m.ctrl.Call(_m, "Query", _param0)
	ret0, _ := ret[0].(*skydb.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query
func (_mr *MockDatabaseMockRecorder) Query(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Query", arg0)
}

// QueryCount mocks base method
func (_m *MockDatabase) QueryCount(_param0 *skydb.Query) (uint64, error) {
	ret := _m.ctrl.Call(_m, "QueryCount", _param0)
	ret0, _ := ret[0].(uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryCount indicates an expected call of QueryCount
func (_mr *MockDatabaseMockRecorder) QueryCount(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "QueryCount", arg0)
}

// RemoteColumnTypes mocks base method
func (_m *MockDatabase) RemoteColumnTypes(_param0 string) (skydb.RecordSchema, error) {
	ret := _m.ctrl.Call(_m, "RemoteColumnTypes", _param0)
	ret0, _ := ret[0].(skydb.RecordSchema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RemoteColumnTypes indicates an expected call of RemoteColumnTypes
func (_mr *MockDatabaseMockRecorder) RemoteColumnTypes(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RemoteColumnTypes", arg0)
}

// RenameSchema mocks base method
func (_m *MockDatabase) RenameSchema(_param0 string, _param1 string, _param2 string) error {
	ret := _m.ctrl.Call(_m, "RenameSchema", _param0, _param1, _param2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RenameSchema indicates an expected call of RenameSchema
func (_mr *MockDatabaseMockRecorder) RenameSchema(arg0, arg1, arg2 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "RenameSchema", arg0, arg1, arg2)
}

// Save mocks base method
func (_m *MockDatabase) Save(_param0 *skydb.Record) error {
	ret := _m.ctrl.Call(_m, "Save", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save
func (_mr *MockDatabaseMockRecorder) Save(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Save", arg0)
}

// SaveSubscription mocks base method
func (_m *MockDatabase) SaveSubscription(_param0 *skydb.Subscription) error {
	ret := _m.ctrl.Call(_m, "SaveSubscription", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSubscription indicates an expected call of SaveSubscription
func (_mr *MockDatabaseMockRecorder) SaveSubscription(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SaveSubscription", arg0)
}

// TableName mocks base method
func (_m *MockDatabase) TableName(_param0 string) string {
	ret := _m.ctrl.Call(_m, "TableName", _param0)
	ret0, _ := ret[0].(string)
	return ret0
}

// TableName indicates an expected call of TableName
func (_mr *MockDatabaseMockRecorder) TableName(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "TableName", arg0)
}

// UserRecordType mocks base method
func (_m *MockDatabase) UserRecordType() string {
	ret := _m.ctrl.Call(_m, "UserRecordType")
	ret0, _ := ret[0].(string)
	return ret0
}

// UserRecordType indicates an expected call of UserRecordType
func (_mr *MockDatabaseMockRecorder) UserRecordType() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UserRecordType")
}
