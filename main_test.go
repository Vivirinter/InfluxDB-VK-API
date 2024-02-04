package main

import (
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type MockedMethodCall struct {
	mock.Mock
}

func (m *MockedMethodCall) CallMethod(method string, options Options, config VkConfig) (string, error) {
	args := m.Called(method, options, config)
	return args.String(0), args.Error(1)
}

func TestMethodCall(t *testing.T) {
	mockObj := new(MockedMethodCall)
	mockMethod := "groups.getById"
	mockOptions := Options{
		"group_id": {"12345"},
	}
	mockConfig := VkConfig{
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
