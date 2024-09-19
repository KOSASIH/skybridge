package network

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNetwork(t *testing.T) {
	n := NewNetwork()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	}))
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	body, err := n.Get(ctx, ts.URL)
	if err != nil {
		t.Errorf("Failed to get URL: %v", err)
	}

	if string(body) != "Hello, world!" {
		t.Errorf("Expected response body to be 'Hello, world!', got %s", body)
	}
}
