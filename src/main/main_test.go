package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Errorf("Failed to create request: %v", err)
	}

	w := httptest.NewRecorder()
	mainHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", w.Code)
	}

	if w.Body.String() != "Welcome to the main application!" {
		t.Errorf("Expected response body to be 'Welcome to the main application!', got %s", w.Body.String())
	}
}
