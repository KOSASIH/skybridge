package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/skybridge/api/types"
	"github.com/skybridge/lib/errors"
	"github.com/skybridge/lib/logging"
)

func TestGateway(t *testing.T) {
	// Create test gateway
	gateway, err := NewGateway(&Config{
		Address:         "localhost",
		Port:            8080,
		TLSKeyFile:      "testdata/tls.key",
		TLSCertFile:     "testdata/tls.crt",
		WebSocketAddress: "localhost",
		WebSocketPort:   8081,
	})
	if err != nil {
		t.Fatal(err)
	}

	// Start test server
	ts := httptest.NewServer(gateway.Router)
	defer ts.Close()

	// Test GET /api/v1/nodes
	req, err := http.NewRequest("GET", ts.URL+"/api/v1/nodes", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+generateToken())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test POST /api/v1/nodes
	req, err = http.NewRequest("POST", ts.URL+"/api/v1/nodes", strings.NewReader(`{"id": "node-1", "name": "Node 1"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+generateToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	// Test GET /api/v1/nodes/{id}
	req, err = http.NewRequest("GET", ts.URL+"/api/v1/nodes/node-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+generateToken())
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test PUT /api/v1/nodes/{id}
	req, err = http.NewRequest("PUT", ts.URL+"/api/v1/nodes/node-1", strings.NewReader(`{"name": "Node 1 Updated"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+generateToken())
	req.Header.Set("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// Test DELETE /api/v1/nodes/{id}
	req, err = http.NewRequest("DELETE", ts.URL+"/api/v1/nodes/node-1", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+generateToken())
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
}

func generateToken() string {
	claims := types.Claims{
		Issuer:    "skybridge",
		Audience:  "skybridge",
		ExpiresAt: time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token SignedString([]byte("secret"))
}
