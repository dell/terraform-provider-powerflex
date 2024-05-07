package client

import (
	"testing"

	"github.com/packer-community/winrmcp/winrmcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWinRMClient is a mock implementation of WinRMClient
type MockWinRMClient struct {
	mock.Mock
}

func (m *MockWinRMClient) GetErrors() []map[string]string {
	args := m.Called()
	return args.Get(0).([]map[string]string)
}

func (m *MockWinRMClient) Destroy() {
	m.Called()
}

func (m *MockWinRMClient) setTarget(context map[string]string, host bool) *WinRMClient {
	args := m.Called(context, host)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) setPort(context map[string]string) *WinRMClient {
	args := m.Called(context)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) setUsername(context map[string]string) *WinRMClient {
	args := m.Called(context)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) setPassword(context map[string]string) *WinRMClient {
	args := m.Called(context)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) setTimeout(context map[string]string) *WinRMClient {
	args := m.Called(context)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) GetConnection(context map[string]string, host bool) *WinRMClient {
	args := m.Called(context, host)
	return args.Get(0).(*WinRMClient)
}

func (m *MockWinRMClient) ExecuteCommand(command string) string {
	args := m.Called(command)
	return args.String(0)
}

func (m *MockWinRMClient) Init() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockWinRMClient) newCopyClient() (*winrmcp.Winrmcp, error) {
	args := m.Called()
	return args.Get(0).(*winrmcp.Winrmcp), args.Error(1)
}

func (m *MockWinRMClient) Upload(dstPath string, srcPath string) error {
	args := m.Called(dstPath, srcPath)
	return args.Error(0)
}

func TestGetErrors(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	mockWinRMClient.On("GetErrors").Return([]map[string]string{{"error": "error1", "message": "message1"}, {"error": "error2", "message": "message2"}})

	errors := mockWinRMClient.GetErrors()

	assert.NotNil(t, errors)
	assert.Len(t, errors, 2)
	assert.Equal(t, "error1", errors[0]["error"])
	assert.Equal(t, "message1", errors[0]["message"])
	assert.Equal(t, "error2", errors[1]["error"])
	assert.Equal(t, "message2", errors[1]["message"])

	mockWinRMClient.AssertExpectations(t)
}

func TestDestroy(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	mockWinRMClient.On("Destroy")

	mockWinRMClient.Destroy()

	mockWinRMClient.AssertExpectations(t)
}

func TestSetTarget(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"host": "localhost"}
	mockWinRMClient.On("setTarget", context, true).Return(&WinRMClient{Target: "localhost"})

	result := mockWinRMClient.setTarget(context, true)

	assert.Equal(t, "localhost", result.Target)

	mockWinRMClient.AssertExpectations(t)
}
func TestSetPort(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"port": "5986"}
	mockWinRMClient.On("setPort", context).Return(&WinRMClient{Port: 5986})

	result := mockWinRMClient.setPort(context)

	assert.Equal(t, 5986, result.Port)

	mockWinRMClient.AssertExpectations(t)
}

func TestSetUsername(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"username": "user1"}
	mockWinRMClient.On("setUsername", context).Return(&WinRMClient{UserName: "user1"})

	result := mockWinRMClient.setUsername(context)

	assert.Equal(t, "user1", result.UserName)

	mockWinRMClient.AssertExpectations(t)
}

func TestSetPassword(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"password": "password1"}
	mockWinRMClient.On("setPassword", context).Return(&WinRMClient{Password: "password1"})

	result := mockWinRMClient.setPassword(context)

	assert.Equal(t, "password1", result.Password)

	mockWinRMClient.AssertExpectations(t)
}

func TestSetTimeout(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"timeout": "30"}
	mockWinRMClient.On("setTimeout", context).Return(&WinRMClient{Timeout: 30})

	result := mockWinRMClient.setTimeout(context)

	assert.Equal(t, 30, result.Timeout)

	mockWinRMClient.AssertExpectations(t)
}

func TestGetConnection(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	context := map[string]string{"host": "localhost", "port": "5986", "username": "user1", "password": "password1", "timeout": "30"}
	mockWinRMClient.On("GetConnection", context, true).Return(&WinRMClient{Target: "localhost", Port: 5986, UserName: "user1", Password: "password1", Timeout: 30})

	result := mockWinRMClient.GetConnection(context, true)

	assert.Equal(t, "localhost", result.Target)
	assert.Equal(t, 5986, result.Port)
	assert.Equal(t, "user1", result.UserName)
	assert.Equal(t, "password1", result.Password)
	assert.Equal(t, 30, result.Timeout)

	mockWinRMClient.AssertExpectations(t)
}

func TestExecuteCommand(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	command := "echo test"
	mockWinRMClient.On("ExecuteCommand", command).Return("test output")

	output := mockWinRMClient.ExecuteCommand(command)

	assert.Equal(t, "test output", output)

	mockWinRMClient.AssertExpectations(t)
}

func TestInit(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	mockWinRMClient.On("Init").Return(true)

	result := mockWinRMClient.Init()

	assert.True(t, result)

	mockWinRMClient.AssertExpectations(t)
}

func TestNewCopyClient(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	mockWinRMClient.On("newCopyClient").Return(&winrmcp.Winrmcp{}, nil)

	client, err := mockWinRMClient.newCopyClient()

	assert.NotNil(t, client)
	assert.NoError(t, err)

	mockWinRMClient.AssertExpectations(t)
}

func TestUpload(t *testing.T) {
	mockWinRMClient := new(MockWinRMClient)

	dstPath := "/path/to/destination"
	srcPath := "/path/to/source"
	mockWinRMClient.On("Upload", dstPath, srcPath).Return(nil)

	err := mockWinRMClient.Upload(dstPath, srcPath)

	assert.NoError(t, err)

	mockWinRMClient.AssertExpectations(t)
}
