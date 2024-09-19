package errors

import (
	"testing"
)

func TestNewError(t *testing.T) {
	e := NewError(404, "Not Found", "Resource not available")
	if e.Code != 404 {
		t.Errorf("Expected error code to be 404, got %d", e.Code)
	}
	if e.Message != "Not Found" {
		t.Errorf("Expected error message to be 'Not Found', got '%s'", e.Message)
	}
	if e.Details != "Resource not available" {
		t.Errorf("Expected error details to be 'Resource not available', got '%s'", e.Details)
	}
}

func TestError_Error(t *testing.T) {
	e := NewError(500, "Internal Server Error")
	if e.Error() != "Error 500: Internal Server Error" {
		t.Errorf("Expected error string to be 'Error 500: Internal Server Error', got '%s'", e.Error())
	}
}

func TestErrorCode(t *testing.T) {
	e := NewError(401, "Unauthorized")
	if e.ErrorCode() != 401 {
		t.Errorf("Expected error code to be 401, got %d", e.ErrorCode())
	}
}

func TestErrorMessage(t *testing.T) {
	e := NewError(403, "Forbidden")
	if e.ErrorMessage() != "Forbidden" {
		t.Errorf("Expected error message to be 'Forbidden', got '%s'", e.ErrorMessage())
	}
}

func TestErrorDetails(t *testing.T) {
	e := NewError(404, "Not Found", "Resource not available")
	if e.ErrorDetails() != "Resource not available" {
		t.Errorf("Expected error details to be 'Resource not available', got '%s'", e.ErrorDetails())
	}
}

func TestIsError(t *testing.T) {
	e := NewError(500, "Internal Server Error")
	if !IsError(e) {
		t.Errorf("Expected IsError to return true, got false")
	}
}
