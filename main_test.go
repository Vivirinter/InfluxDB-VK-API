package tests

import (
	"github.com/Vivirinter/InfluxDB-VK-API/cmd/main"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type MockedMethodCall struct {
	mock.Mock
}

func (m *MockedMethodCall) CallMethod(method string, options main.Options, config main.VkConfig) (string, error) {
	args := m.Called(method, options, config)
	return args.String(0), args.Error(1)
}

func TestMethodCall(t *testing.T) {
	mockObj := new(MockedMethodCall)
	mockMethod := "groups.getById"
	mockOptions := main.Options{
		"group_id": {"12345"},
	}
	mockConfig := main.VkConfig{
		Token:   "token",
		GroupId: "groupId",
		Version: "5.199",
	}
	mockResponse := "mockResponse"
	mockObj.On("CallMethod", mockMethod, mockOptions, mockConfig).Return(mockResponse, nil)

	result, _ := mockObj.CallMethod(mockMethod, mockOptions, mockConfig)

	if !reflect.DeepEqual(result, mockResponse) {
		t.Fatalf("Expected \"%s\" but got \"%s\"", mockResponse, result)
	}

	mockObj.AssertExpectations(t)
}
